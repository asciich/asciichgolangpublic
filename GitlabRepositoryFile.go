package asciichgolangpublic

import (
	"errors"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	"github.com/asciich/asciichgolangpublic/encoding/base64"
	aerrors "github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
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

func (g *GitlabRepositoryFile) Delete(commitMessage string, verbose bool) (err error) {
	if commitMessage == "" {
		return aerrors.TracedErrorEmptyString("commitMessage")
	}

	nativeClient, projectId, err := g.GetNativeRepositoryFilesClientAndProjectId()
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl()
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
		branchName, err = g.GetDefaultBranchName()
		if err != nil {
			return err
		}
	}

	if branchName == "" {
		return aerrors.TracedError("branchName is empty string after evaluation")
	}

	exits, err := g.Exists()
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
			return aerrors.TracedErrorf(
				"Failed to delete '%s' in branch '%s' on '%s': '%w'",
				fileName,
				branchName,
				projectUrl,
				err,
			)
		}

		if verbose {
			logging.LogChangedf(
				"File '%s' in branch '%s' of gitlab project '%s' deleted.",
				fileName,
				branchName,
				projectUrl,
			)
		}
	} else {
		logging.LogInfof(
			"File '%s' in branch '%s' of gitlab project '%s' is already absent.",
			fileName,
			branchName,
			projectUrl,
		)
	}

	return err
}

func (g *GitlabRepositoryFile) Exists() (fileExists bool, err error) {
	_, err = g.GetNativeRepositoryFile()
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
		return "", aerrors.TracedError("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabRepositoryFile) GetContentAsBytes(verbose bool) (content []byte, err error) {
	content, _, err = g.GetContentAsBytesAndCommitHash(verbose)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (g *GitlabRepositoryFile) GetContentAsBytesAndCommitHash(verbose bool) (content []byte, sha256sum string, err error) {
	nativeRepoFile, err := g.GetNativeRepositoryFile()
	if err != nil {
		return nil, "", err
	}

	contentBase64 := nativeRepoFile.Content

	content, err = base64.DecodeStringAsBytes(contentBase64)
	if err != nil {
		return nil, "", err
	}

	sha256sum = nativeRepoFile.SHA256

	if sha256sum == "" {
		return nil, "", aerrors.TracedError("sha256sum is empty string after evaluation.")
	}

	return content, sha256sum, nil
}

func (g *GitlabRepositoryFile) GetContentAsString(verbose bool) (content string, err error) {
	contentBytes, err := g.GetContentAsBytes(verbose)
	if err != nil {
		return "", err
	}

	content = string(contentBytes)
	return content, nil
}

func (g *GitlabRepositoryFile) GetDefaultBranchName() (defaultBranchName string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	defaultBranchName, err = project.GetDefaultBranchName()
	if err != nil {
		return "", err
	}

	return defaultBranchName, nil
}

func (g *GitlabRepositoryFile) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, aerrors.TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabRepositoryFile) GetNativeRepositoryFile() (nativeFile *gitlab.File, err error) {
	nativeRepositoryFilesClient, err := g.GetNativeRepositoryFilesClient()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId()
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
		branchName, err := g.GetDefaultBranchName()
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
			return nil, aerrors.TracedErrorf("%w, %w", ErrGitlabRepositoryFileDoesNotExist, err)
		}
		return nil, aerrors.TracedErrorf("Unable to get native file: '%w'", err)
	}

	if nativeFile == nil {
		return nil, aerrors.TracedError("nativeFile is nil after evaluation")
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

func (g *GitlabRepositoryFile) GetNativeRepositoryFilesClientAndProjectId() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, projectId int, err error) {
	nativeRepositoryFilesClient, err = g.GetNativeRepositoryFilesClient()
	if err != nil {
		return nil, -1, err
	}

	projectId, err = g.GetProjectId()
	if err != nil {
		return nil, -1, err
	}

	return nativeRepositoryFilesClient, projectId, nil
}

func (g *GitlabRepositoryFile) GetPath() (path string, err error) {
	if g.Path == "" {
		return "", aerrors.TracedErrorf("Path not set")
	}

	return g.Path, nil
}

func (g *GitlabRepositoryFile) GetProjectId() (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	projectId, err = project.GetId()
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabRepositoryFile) GetProjectUrl() (projectUrl string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	projectUrl, err = project.GetProjectUrl()
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

func (g *GitlabRepositoryFile) GetSha256CheckSum() (checkSum string, err error) {
	rawResponse, err := g.GetNativeRepositoryFile()
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
		return "", aerrors.TracedErrorf("SHA256 checksum is empty string after evalutaion for repository file '%s' in branch '%s'.", filePath, branchName)
	}

	return checkSum, nil
}

func (g *GitlabRepositoryFile) IsBranchNameSet() (isSet bool) {
	return g.BranchName != ""
}

func (g *GitlabRepositoryFile) MustDelete(commitMessage string, verbose bool) {
	err := g.Delete(commitMessage, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFile) MustExists() (fileExists bool) {
	fileExists, err := g.Exists()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileExists
}

func (g *GitlabRepositoryFile) MustGetBranchName() (branchName string) {
	branchName, err := g.GetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchName
}

func (g *GitlabRepositoryFile) MustGetContentAsBytes(verbose bool) (content []byte) {
	content, err := g.GetContentAsBytes(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (g *GitlabRepositoryFile) MustGetContentAsBytesAndCommitHash(verbose bool) (content []byte, sha256sum string) {
	content, sha256sum, err := g.GetContentAsBytesAndCommitHash(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content, sha256sum
}

func (g *GitlabRepositoryFile) MustGetContentAsString(verbose bool) (content string) {
	content, err := g.GetContentAsString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (g *GitlabRepositoryFile) MustGetDefaultBranchName() (defaultBranchName string) {
	defaultBranchName, err := g.GetDefaultBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return defaultBranchName
}

func (g *GitlabRepositoryFile) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabRepositoryFile) MustGetNativeRepositoryFile() (nativeFile *gitlab.File) {
	nativeFile, err := g.GetNativeRepositoryFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeFile
}

func (g *GitlabRepositoryFile) MustGetNativeRepositoryFilesClient() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService) {
	nativeRepositoryFilesClient, err := g.GetNativeRepositoryFilesClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeRepositoryFilesClient
}

func (g *GitlabRepositoryFile) MustGetNativeRepositoryFilesClientAndProjectId() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, projectId int) {
	nativeRepositoryFilesClient, projectId, err := g.GetNativeRepositoryFilesClientAndProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeRepositoryFilesClient, projectId
}

func (g *GitlabRepositoryFile) MustGetPath() (path string) {
	path, err := g.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (g *GitlabRepositoryFile) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabRepositoryFile) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabRepositoryFile) MustGetRepositoryFiles() (repositoryFiles *GitlabRepositoryFiles) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repositoryFiles
}

func (g *GitlabRepositoryFile) MustGetSha256CheckSum() (checkSum string) {
	checkSum, err := g.GetSha256CheckSum()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return checkSum
}

func (g *GitlabRepositoryFile) MustSetBranchName(branchName string) {
	err := g.SetBranchName(branchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFile) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFile) MustSetPath(path string) {
	err := g.SetPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFile) MustWriteFileContentByBytes(content []byte, commitMessage string, verbose bool) {
	err := g.WriteFileContentByBytes(content, commitMessage, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFile) MustWriteFileContentByString(content string, commitMessage string, verbose bool) {
	err := g.WriteFileContentByString(content, commitMessage, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabRepositoryFile) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return aerrors.TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabRepositoryFile) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return aerrors.TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabRepositoryFile) SetPath(path string) (err error) {
	if path == "" {
		return aerrors.TracedErrorf("path is empty string")
	}

	g.Path = path

	return nil
}

func (g *GitlabRepositoryFile) WriteFileContentByBytes(content []byte, commitMessage string, verbose bool) (err error) {
	if content == nil {
		return aerrors.TracedErrorNil("content")
	}

	if commitMessage == "" {
		return aerrors.TracedErrorEmptyString("commitMessage")
	}

	exists, err := g.Exists()
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

	projectId, err := g.GetProjectId()
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl()
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
		branchName, err = g.GetDefaultBranchName()
		if err != nil {
			return err
		}
	}
	if branchName == "" {
		return aerrors.TracedError("Internal error: branchName is empty string after evaluation.")
	}

	contentString := string(content)

	if exists {
		currentContent, err := g.GetContentAsBytes(verbose)
		if err != nil {
			return err
		}

		if aslices.ByteSlicesEqual(currentContent, content) {
			if verbose {
				logging.LogInfof(
					"Content of Gitlab repository file '%s' in project '%s' is already up to date.",
					fileName,
					projectUrl,
				)
			}
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
				return aerrors.TracedErrorf("Unable to update file: '%w'", err)
			}

			if verbose {
				logging.LogChangedf(
					"Content of Gitlab repository file '%s' in project '%s' updated.",
					fileName,
					projectUrl,
				)
			}
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
			return aerrors.TracedErrorf("Unable to create file in gitlab project: '%w'", err)
		}

		if verbose {
			logging.LogChangedf("Created file '%s' in Gitlab project '%s'.", fileName, projectUrl)
		}
	}

	return nil
}

func (g *GitlabRepositoryFile) WriteFileContentByString(content string, commitMessage string, verbose bool) (err error) {
	err = g.WriteFileContentByBytes([]byte(content), commitMessage, verbose)
	if err != nil {
		return err
	}

	return nil
}
