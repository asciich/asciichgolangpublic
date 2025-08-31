package asciichgolangpublic

import (
	"context"
	"errors"
	"slices"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/pkg/encodingutils/base64utils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var ErrGitlabRepositoryFileDoesNotExist = errors.New("Gitlab repository file does not exist")

type GitlabRepositoryFile struct {
	gitlabProject *GitlabProject
	Path          string
	BranchName    string
}

func NewGitlabRepositoryFile() (g *GitlabRepositoryFile) {
	return new(GitlabRepositoryFile)
}

func (g *GitlabRepositoryFile) Delete(ctx context.Context, commitMessage string) (err error) {
	if commitMessage == "" {
		return tracederrors.TracedErrorEmptyString("commitMessage")
	}

	nativeClient, projectId, err := g.GetNativeRepositoryFilesClientAndProjectId(ctx)
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl(ctx)
	if err != nil {
		return err
	}

	fileName, err := g.GetPath()
	if err != nil {
		return err
	}

	var branchName = ""
	if g.IsBranchNameSet() {
		branchName, err = g.GetBranchName()
		if err != nil {
			return err
		}
	} else {
		branchName, err = g.GetDefaultBranchName(ctx)
		if err != nil {
			return err
		}
	}

	if branchName == "" {
		return tracederrors.TracedError("branchName is empty string after evaluation")
	}

	exits, err := g.Exists(ctx)
	if err != nil {
		return err
	}

	if exits {
		_, err = nativeClient.DeleteFile(
			projectId,
			fileName,
			&gitlab.DeleteFileOptions{
				Branch:        &branchName,
				CommitMessage: &commitMessage,
			},
		)
		if err != nil {
			return tracederrors.TracedErrorf(
				"Failed to delete '%s' in branch '%s' on '%s': '%w'",
				fileName,
				branchName,
				projectUrl,
				err,
			)
		}

		logging.LogChangedByCtxf(ctx, "File '%s' in branch '%s' of gitlab project '%s' deleted.", fileName, branchName, projectUrl)
	} else {
		logging.LogInfoByCtxf(ctx, "File '%s' in branch '%s' of gitlab project '%s' is already absent.", fileName, branchName, projectUrl)
	}

	return err
}

func (g *GitlabRepositoryFile) Exists(ctx context.Context) (fileExists bool, err error) {
	_, err = g.GetNativeRepositoryFile(ctx)
	if err != nil {
		if errors.Is(err, ErrGitlabRepositoryFileDoesNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (g *GitlabRepositoryFile) GetBranchName() (branchName string, err error) {
	if g.BranchName == "" {
		return "", tracederrors.TracedError("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabRepositoryFile) GetContentAsBytes(ctx context.Context) (content []byte, err error) {
	content, _, err = g.GetContentAsBytesAndCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (g *GitlabRepositoryFile) GetContentAsBytesAndCommitHash(ctx context.Context) (content []byte, sha256sum string, err error) {
	nativeRepoFile, err := g.GetNativeRepositoryFile(ctx)
	if err != nil {
		return nil, "", err
	}

	contentBase64 := nativeRepoFile.Content

	content, err = base64utils.DecodeStringAsBytes(contentBase64)
	if err != nil {
		return nil, "", err
	}

	sha256sum = nativeRepoFile.SHA256

	if sha256sum == "" {
		return nil, "", tracederrors.TracedError("sha256sum is empty string after evaluation.")
	}

	return content, sha256sum, nil
}

func (g *GitlabRepositoryFile) GetContentAsString(ctx context.Context) (content string, err error) {
	contentBytes, err := g.GetContentAsBytes(ctx)
	if err != nil {
		return "", err
	}

	content = string(contentBytes)
	return content, nil
}

func (g *GitlabRepositoryFile) GetDefaultBranchName(ctx context.Context) (defaultBranchName string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	defaultBranchName, err = project.GetDefaultBranchName(ctx)
	if err != nil {
		return "", err
	}

	return defaultBranchName, nil
}

func (g *GitlabRepositoryFile) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, tracederrors.TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabRepositoryFile) GetNativeRepositoryFile(ctx context.Context) (nativeFile *gitlab.File, err error) {
	nativeRepositoryFilesClient, err := g.GetNativeRepositoryFilesClient()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return nil, err
	}

	fileName, err := g.GetPath()
	if err != nil {
		return nil, err
	}

	var getFileOptions *gitlab.GetFileOptions
	if g.IsBranchNameSet() {
		branchName, err := g.GetBranchName()
		if err != nil {
			return nil, err
		}

		getFileOptions = new(gitlab.GetFileOptions)
		getFileOptions.Ref = &branchName
	} else {
		branchName, err := g.GetDefaultBranchName(ctx)
		if err != nil {
			return nil, err
		}

		getFileOptions = new(gitlab.GetFileOptions)
		getFileOptions.Ref = &branchName
	}

	nativeFile, _, err = nativeRepositoryFilesClient.GetFile(
		projectId,
		fileName,
		getFileOptions,
	)
	if err != nil {
		if err.Error() == "404 Not Found" {
			return nil, tracederrors.TracedErrorf("%w, %w", ErrGitlabRepositoryFileDoesNotExist, err)
		}
		return nil, tracederrors.TracedErrorf("Unable to get native file: '%w'", err)
	}

	if nativeFile == nil {
		return nil, tracederrors.TracedError("nativeFile is nil after evaluation")
	}

	return nativeFile, nil
}

func (g *GitlabRepositoryFile) GetNativeRepositoryFilesClient() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	nativeRepositoryFilesClient, err = repositoryFiles.GetNativeRepositoryFilesClient()
	if err != nil {
		return nil, err
	}

	return nativeRepositoryFilesClient, nil
}

func (g *GitlabRepositoryFile) GetNativeRepositoryFilesClientAndProjectId(ctx context.Context) (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, projectId int, err error) {
	nativeRepositoryFilesClient, err = g.GetNativeRepositoryFilesClient()
	if err != nil {
		return nil, -1, err
	}

	projectId, err = g.GetProjectId(ctx)
	if err != nil {
		return nil, -1, err
	}

	return nativeRepositoryFilesClient, projectId, nil
}

func (g *GitlabRepositoryFile) GetPath() (path string, err error) {
	if g.Path == "" {
		return "", tracederrors.TracedErrorf("Path not set")
	}

	return g.Path, nil
}

func (g *GitlabRepositoryFile) GetProjectId(ctx context.Context) (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	projectId, err = project.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabRepositoryFile) GetProjectUrl(ctx context.Context) (projectUrl string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	projectUrl, err = project.GetProjectUrl(ctx)
	if err != nil {
		return "", err
	}

	return projectUrl, nil
}

func (g *GitlabRepositoryFile) GetRepositoryFiles() (repositoryFiles *GitlabRepositoryFiles, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	repositoryFiles, err = gitlabProject.GetRepositoryFiles()
	if err != nil {
		return nil, err
	}

	return repositoryFiles, nil
}

func (g *GitlabRepositoryFile) GetSha256CheckSum(ctx context.Context) (checkSum string, err error) {
	rawResponse, err := g.GetNativeRepositoryFile(ctx)
	if err != nil {
		return "", err
	}

	filePath, err := g.GetPath()
	if err != nil {
		return "", err
	}

	branchName, err := g.GetBranchName()
	if err != nil {
		return "", err
	}

	checkSum = rawResponse.SHA256

	if checkSum == "" {
		return "", tracederrors.TracedErrorf("SHA256 checksum is empty string after evalutaion for repository file '%s' in branch '%s'.", filePath, branchName)
	}

	return checkSum, nil
}

func (g *GitlabRepositoryFile) IsBranchNameSet() (isSet bool) {
	return g.BranchName != ""
}

func (g *GitlabRepositoryFile) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return tracederrors.TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabRepositoryFile) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabRepositoryFile) SetPath(path string) (err error) {
	if path == "" {
		return tracederrors.TracedErrorf("path is empty string")
	}

	g.Path = path

	return nil
}

func (g *GitlabRepositoryFile) WriteFileContentByBytes(ctx context.Context, content []byte, commitMessage string) (err error) {
	if content == nil {
		return tracederrors.TracedErrorNil("content")
	}

	if commitMessage == "" {
		return tracederrors.TracedErrorEmptyString("commitMessage")
	}

	exists, err := g.Exists(ctx)
	if err != nil {
		return err
	}

	nativeClient, err := g.GetNativeRepositoryFilesClient()
	if err != nil {
		return err
	}

	fileName, err := g.GetPath()
	if err != nil {
		return err
	}

	projectId, err := g.GetProjectId(ctx)
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl(ctx)
	if err != nil {
		return err
	}

	branchName := ""
	if g.IsBranchNameSet() {
		branchName, err = g.GetBranchName()
		if err != nil {
			return err
		}
	} else {
		branchName, err = g.GetDefaultBranchName(ctx)
		if err != nil {
			return err
		}
	}
	if branchName == "" {
		return tracederrors.TracedError("Internal error: branchName is empty string after evaluation.")
	}

	contentString := string(content)

	if exists {
		currentContent, err := g.GetContentAsBytes(ctx)
		if err != nil {
			return err
		}

		if slices.Equal(currentContent, content) {
			logging.LogInfoByCtxf(ctx, "Content of Gitlab repository file '%s' in project '%s' is already up to date.", fileName, projectUrl)
		} else {
			updateOptions := new(gitlab.UpdateFileOptions)
			updateOptions.CommitMessage = &commitMessage
			updateOptions.Branch = &branchName
			updateOptions.Content = &contentString

			_, _, err := nativeClient.UpdateFile(
				projectId,
				fileName,
				updateOptions,
				nil,
			)
			if err != nil {
				return tracederrors.TracedErrorf("Unable to update file: '%w'", err)
			}

			logging.LogChangedByCtxf(ctx, "Content of Gitlab repository file '%s' in project '%s' updated.", fileName, projectUrl)
		}
	} else {
		createOptions := new(gitlab.CreateFileOptions)

		createOptions.Content = &contentString

		createOptions.CommitMessage = &commitMessage
		createOptions.Branch = &branchName

		_, _, err := nativeClient.CreateFile(
			projectId,
			fileName,
			createOptions,
			nil,
		)
		if err != nil {
			return tracederrors.TracedErrorf("Unable to create file in gitlab project: '%w'", err)
		}

		logging.LogChangedByCtxf(ctx, "Created file '%s' in Gitlab project '%s'.", fileName, projectUrl)
	}

	return nil
}

func (g *GitlabRepositoryFile) WriteFileContentByString(ctx context.Context, content string, commitMessage string) (err error) {
	err = g.WriteFileContentByBytes(ctx, []byte(content), commitMessage)
	if err != nil {
		return err
	}

	return nil
}
