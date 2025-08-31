package asciichgolangpublic

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabBranch struct {
	gitlabProject *GitlabProject
	name          string
}

func NewGitlabBranch() (g *GitlabBranch) {
	return new(GitlabBranch)
}

func (g *GitlabBranch) CopyFileToBranch(ctx context.Context, filePath string, targetBranch *GitlabBranch) (targetFile *GitlabRepositoryFile, err error) {
	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
	}

	if targetBranch == nil {
		return nil, tracederrors.TracedErrorNil("targetBranch")
	}

	sourceFile, err := g.GetRepositoryFile(ctx, filePath)
	if err != nil {
		return nil, err
	}

	targetFile, err = targetBranch.GetRepositoryFile(ctx, filePath)
	if err != nil {
		return nil, err
	}

	sourceSha, err := sourceFile.GetSha256CheckSum(ctx)
	if err != nil {
		return nil, err
	}

	destSha := ""

	exists, err := targetFile.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if exists {
		destSha, err = targetFile.GetSha256CheckSum(ctx)
		if err != nil {
			return nil, err
		}
	}

	sourceBranchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	targetBranchName, err := targetBranch.GetName()
	if err != nil {
		return nil, err
	}

	if sourceSha == destSha {
		logging.LogInfoByCtxf(ctx, "File '%s' is equal in source branch '%s' and target branch '%s' have already equal content. Skip copy.", filePath, sourceBranchName, targetBranchName)
	} else {
		content, commitHash, err := sourceFile.GetContentAsBytesAndCommitHash(ctx)
		if err != nil {
			return nil, err
		}

		commitMessage := fmt.Sprintf(
			"Copy '%s' from commit '%s' of branch '%s' to target branch '%s'.",
			filePath,
			commitHash,
			sourceBranchName,
			targetBranchName,
		)

		err = targetFile.WriteFileContentByBytes(ctx, content, commitMessage)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "File '%s' is copied from commit '%s' of source branch '%s' to target branch '%s'.", filePath, commitHash, sourceBranchName, targetBranchName)
	}

	return targetFile, nil
}

func (g *GitlabBranch) CreateFromDefaultBranch(ctx context.Context) (err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return err
	}

	branchName, err := g.GetName()
	if err != nil {
		return err
	}

	_, err = branches.CreateBranchFromDefaultBranch(ctx, branchName)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabBranch) CreateMergeRequest(ctx context.Context, options *GitlabCreateMergeRequestOptions) (mergeRequest *GitlabMergeRequest, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	mergeRequests, err := g.GetMergeRequests()
	if err != nil {
		return nil, err
	}

	branchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	optionsToUse := options.GetDeepCopy()
	err = optionsToUse.SetSourceBranchName(branchName)
	if err != nil {
		return nil, err
	}

	mergeRequest, err = mergeRequests.CreateMergeRequest(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabBranch) Delete(ctx context.Context, options *GitlabDeleteBranchOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	nativeClient, projectId, err := g.GetNativeBranchesClientAndId(ctx)
	if err != nil {
		return err
	}

	branchName, err := g.GetName()
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl(ctx)
	if err != nil {
		return err
	}

	exists, err := g.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = nativeClient.DeleteBranch(
			projectId,
			branchName,
		)
		if err != nil {
			return tracederrors.TracedErrorf(
				"Delete branch '%s' in gitlab project %s failed: %w",
				branchName,
				projectUrl,
				err,
			)
		}

		if options.SkipWaitForDeletion {
			exists = false
		} else {
			// Deleting is not instantaneous so lets check if deleted branch is really absent.
			// Especially the branch listing on gitlab side tends to race conditions.
			for i := 0; i < 30; i++ {
				gitlabBranches, err := g.GetGitlabBranches()
				if err != nil {
					return err
				}

				branchNames, err := gitlabBranches.GetBranchNames(contextutils.WithSilent(ctx))
				if err != nil {
					return err
				}

				exists = slices.Contains(branchNames, branchName)
				if exists {
					time.Sleep(1000 * time.Millisecond)
					logging.LogInfoByCtxf(ctx, "Wait for branch '%s' to be deleted in %s .", branchName, projectUrl)
				} else {
					break
				}
			}
		}

		if exists {
			return tracederrors.TracedErrorf("Internal error: failed to delete '%s' in %s", branchName, projectUrl)
		}

		logging.LogChangedByCtxf(ctx, "Deleted branch '%s' in gitlab project %s .", branchName, projectUrl)
	} else {
		logging.LogInfoByCtxf(ctx, "Branch '%s' is already absent on %s .", branchName, projectUrl)
	}

	return nil
}

func (g *GitlabBranch) DeleteRepositoryFile(ctx context.Context, filePath string, commitMessage string) (err error) {
	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	if commitMessage == "" {
		return tracederrors.TracedErrorEmptyString("commitMessage")
	}

	fileToDelete, err := g.GetRepositoryFile(ctx, filePath)
	if err != nil {
		return err
	}

	err = fileToDelete.Delete(ctx, commitMessage)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabBranch) Exists(ctx context.Context) (exists bool, err error) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId(ctx)
	if err != nil {
		return false, err
	}

	branchName, err := g.GetName()
	if err != nil {
		return false, err
	}

	_, _, err = nativeClient.GetBranch(
		projectId,
		branchName,
		nil,
	)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}

		return false, tracederrors.TracedErrorf("Failed to evaluate if branch exists: '%w'", err)
	}

	return true, nil
}

func (g *GitlabBranch) GetBranches() (branches *GitlabBranches, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branches, err = project.GetBranches()
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (g *GitlabBranch) GetDeepCopy() (copy *GitlabBranch) {
	copy = NewGitlabBranch()

	*copy = *g

	if g.gitlabProject != nil {
		copy.gitlabProject = g.gitlabProject.GetDeepCopy()
	}

	return copy
}

func (g *GitlabBranch) GetGitlab() (gitlab *GitlabInstance, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlab, err = project.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func (g *GitlabBranch) GetGitlabBranches() (branches *GitlabBranches, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branches, err = gitlabProject.GetBranches()
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (g *GitlabBranch) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, tracederrors.TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabBranch) GetLatestCommit(ctx context.Context) (latestCommit *GitlabCommit, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	latestCommit, err = gitlabProject.GetLatestCommit(ctx, branchName)
	if err != nil {
		return nil, err
	}

	return latestCommit, err
}

func (g *GitlabBranch) GetLatestCommitHashAsString(ctx context.Context) (commitHash string, err error) {
	rawResponse, err := g.GetRawResponse(ctx)
	if err != nil {
		return "", err
	}

	rawCommit := rawResponse.Commit
	if rawCommit == nil {
		return "", tracederrors.TracedError("rawCommit is nil in get latest commit hash as string")
	}

	commitHash = rawCommit.ID
	if commitHash == "" {
		return "", tracederrors.TracedError("commitHash is empty string after evaluation")
	}

	branchName, err := g.GetName()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Latest commit of branch '%s' is '%s'", branchName, commitHash)

	return commitHash, nil
}

func (g *GitlabBranch) GetMergeRequests() (mergeRequests *GitlabProjectMergeRequests, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	mergeRequests, err = project.GetMergeRequests()
	if err != nil {
		return nil, err
	}

	return mergeRequests, nil
}

func (g *GitlabBranch) GetName() (name string, err error) {
	if g.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GitlabBranch) GetNativeBranchesClient() (nativeClient *gitlab.BranchesService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeBranchesClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (g *GitlabBranch) GetNativeBranchesClientAndId(ctx context.Context) (nativeClient *gitlab.BranchesService, projectId int, err error) {
	nativeClient, err = g.GetNativeBranchesClient()
	if err != nil {
		return nil, -1, err
	}

	projectId, err = g.GetProjectId(ctx)
	if err != nil {
		return nil, -1, err
	}

	return nativeClient, projectId, nil
}

func (g *GitlabBranch) GetProjectId(ctx context.Context) (projectId int, err error) {
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

func (g *GitlabBranch) GetProjectUrl(ctx context.Context) (projectUrl string, err error) {
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

func (g *GitlabBranch) GetRawResponse(ctx context.Context) (rawResponse *gitlab.Branch, err error) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId(ctx)
	if err != nil {
		return nil, err
	}

	branchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	rawResponse, _, err = nativeClient.GetBranch(
		projectId,
		branchName,
		nil,
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to get branch: '%w'", err)
	}

	if rawResponse == nil {
		return nil, tracederrors.TracedError("rawResponse for GitlabBranch is nil after evaluation")
	}

	return rawResponse, nil
}

func (g *GitlabBranch) GetRepositoryFile(ctx context.Context, filePath string) (repositoryFile *GitlabRepositoryFile, err error) {
	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	repositoryFile, err = gitlabProject.GetRepositoryFile(
		ctx,
		&GitlabGetRepositoryFileOptions{
			BranchName: branchName,
			Path:       filePath,
		},
	)
	if err != nil {
		return nil, err
	}

	return repositoryFile, nil
}

func (g *GitlabBranch) GetRepositoryFileSha256Sum(ctx context.Context, filePath string) (sha256sum string, err error) {
	if filePath == "" {
		return "", tracederrors.TracedErrorEmptyString("filePath")
	}

	repostioryFile, err := g.GetRepositoryFile(ctx, filePath)
	if err != nil {
		return "", err
	}

	sha256sum, err = repostioryFile.GetSha256CheckSum(ctx)
	if err != nil {
		return "", err
	}

	return sha256sum, nil
}

func (g *GitlabBranch) ReadFileContentAsString(ctx context.Context, options *GitlabReadFileOptions) (content string, err error) {
	if options == nil {
		return "", tracederrors.TracedErrorNil("options")
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	optionsToUse := options.GetDeepCopy()

	branchName, err := g.GetName()
	if err != nil {
		return "", err
	}

	err = optionsToUse.SetBranchName(branchName)
	if err != nil {
		return "", err
	}

	content, err = gitlabProject.ReadFileContentAsString(ctx, optionsToUse)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (g *GitlabBranch) RepositoryFileExists(ctx context.Context, filePath string) (exists bool, err error) {
	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	repositoryFile, err := g.GetRepositoryFile(ctx, filePath)
	if err != nil {
		return false, err
	}

	exists, err = repositoryFile.Exists(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (g *GitlabBranch) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabBranch) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}

func (g *GitlabBranch) SyncFilesToBranch(ctx context.Context, options *GitlabSyncBranchOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	targetBranchName, err := options.GetTargetBranchName()
	if err != nil {
		return err
	}

	sourceBranchName, err := g.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Sync files from source branch '%s' to target branch '%s' started.", sourceBranchName, targetBranchName)

	branches, err := g.GetBranches()
	if err != nil {
		return err
	}

	targetBranch, err := branches.GetBranchByName(targetBranchName)
	if err != nil {
		return err
	}

	pathsToSync, err := options.GetPathsToSync()
	if err != nil {
		return err
	}

	filesToSync, err := branches.GetFilesFromListWithDiffBetweenBranches(ctx, g, targetBranch, pathsToSync)
	if err != nil {
		return err
	}

	if len(filesToSync) <= 0 {
		logging.LogInfoByCtxf(ctx, "All '%d' files to sync from branch '%s' to target branch '%s' are already up to date.", len(pathsToSync), sourceBranchName, targetBranchName)
	} else {
		for _, pathToSync := range filesToSync {
			_, err = g.CopyFileToBranch(ctx, pathToSync, targetBranch)
			if err != nil {
				return err
			}
		}

		logging.LogChangedByCtxf(ctx, "Synced '%d' files to sync from branch '%s' to target branch '%s'. '%d' files were already up to date so there was no need for a sync.", len(filesToSync), sourceBranchName, targetBranchName, len(pathsToSync)-len(filesToSync))
	}

	logging.LogInfoByCtxf(ctx, "Sync files from source branch '%s' to target branch '%s' finished.", sourceBranchName, targetBranchName)

	return nil
}

func (g *GitlabBranch) SyncFilesToBranchUsingMergeRequest(ctx context.Context, options *GitlabSyncBranchOptions) (createdMergeRequest *GitlabMergeRequest, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	targetBranchName, err := options.GetTargetBranchName()
	if err != nil {
		return nil, err
	}

	sourceBranchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	branches, err := g.GetBranches()
	if err != nil {
		return nil, err
	}

	targetBranch, err := branches.GetBranchByName(targetBranchName)
	if err != nil {
		return nil, err
	}

	pathsToSync, err := options.GetPathsToSync()
	if err != nil {
		return nil, err
	}

	filesToSync, err := branches.GetFilesFromListWithDiffBetweenBranches(ctx, g, targetBranch, pathsToSync)
	if err != nil {
		return nil, err
	}

	syncBranchName := fmt.Sprintf(
		"sync-files-%s-to-%s",
		sourceBranchName,
		targetBranchName,
	)

	if len(filesToSync) > 0 {

		optionsToUse := options.GetDeepCopy()

		err = optionsToUse.SetTargetBranchNameAndUnsetTargetBranchObject(syncBranchName)
		if err != nil {
			return nil, err
		}

		syncBranch, err := branches.CreateBranch(
			ctx,
			&GitlabCreateBranchOptions{
				SourceBranchName:    targetBranchName,
				BranchName:          syncBranchName,
				FailIfAlreadyExists: true,
			},
		)
		if err != nil {
			return nil, err
		}

		err = g.SyncFilesToBranch(ctx, optionsToUse)
		if err != nil {
			return nil, err
		}

		courceCommitHash, err := g.GetLatestCommitHashAsString(ctx)
		if err != nil {
			return nil, err
		}

		mrTitle := fmt.Sprintf(
			"Sync files from '%s' hash '%s' to '%s'",
			sourceBranchName,
			courceCommitHash,
			targetBranchName,
		)

		mrDescription := fmt.Sprintf(
			"Sync files from branch '%s' hash '%s' to '%s'.",
			sourceBranchName,
			courceCommitHash,
			targetBranchName,
		)

		createdMergeRequest, err = syncBranch.CreateMergeRequest(
			ctx,
			&GitlabCreateMergeRequestOptions{
				TargetBranchName:                targetBranchName,
				Title:                           mrTitle,
				Description:                     mrDescription,
				SquashEnabled:                   true,
				DeleteSourceBranchOnMerge:       true,
				FailIfMergeRequestAlreadyExists: true,
				AssignToSelf:                    true,
			},
		)
		if err != nil {
			return nil, err
		}

		mrUrl, err := createdMergeRequest.GetUrlAsString(ctx)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Created merge request %s to sync '%d' files from branch '%s' to target branch '%s'.", mrUrl, len(filesToSync), sourceBranchName, targetBranchName)
	} else {
		logging.LogInfof(
			"All '%d' files are in sync between branch '%s' and '%s'. No merge request was created.",
			len(pathsToSync),
			syncBranchName,
			targetBranchName,
		)
	}

	return createdMergeRequest, nil
}

func (g *GitlabBranch) WriteFileContent(ctx context.Context, options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	optionsToUse := options.GetDeepCopy()
	optionsToUse.BranchName = branchName

	gitlabRepositoryFile, err = gitlabProject.WriteFileContent(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	return gitlabRepositoryFile, nil
}
