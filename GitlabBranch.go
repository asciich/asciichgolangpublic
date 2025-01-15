package asciichgolangpublic

import (
	"fmt"
	"strings"
	"time"

	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabBranch struct {
	gitlabProject *GitlabProject
	name          string
}

func NewGitlabBranch() (g *GitlabBranch) {
	return new(GitlabBranch)
}

func (g *GitlabBranch) CopyFileToBranch(filePath string, targetBranch *GitlabBranch, verbose bool) (targetFile *GitlabRepositoryFile, err error) {
	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
	}

	if targetBranch == nil {
		return nil, tracederrors.TracedErrorNil("targetBranch")
	}

	sourceFile, err := g.GetRepositoryFile(filePath, verbose)
	if err != nil {
		return nil, err
	}

	targetFile, err = targetBranch.GetRepositoryFile(filePath, verbose)
	if err != nil {
		return nil, err
	}

	sourceSha, err := sourceFile.GetSha256CheckSum()
	if err != nil {
		return nil, err
	}

	destSha := ""

	exists, err := targetFile.Exists()
	if err != nil {
		return nil, err
	}

	if exists {
		destSha, err = targetFile.GetSha256CheckSum()
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
		if verbose {
			logging.LogInfof(
				"File '%s' is equal in source branch '%s' and target branch '%s' have already equal content. Skip copy.",
				filePath,
				sourceBranchName,
				targetBranchName,
			)
		}
	} else {
		content, commitHash, err := sourceFile.GetContentAsBytesAndCommitHash(verbose)
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

		err = targetFile.WriteFileContentByBytes(content, commitMessage, verbose)
		if err != nil {
			return nil, err
		}

		if verbose {
			logging.LogChangedf(
				"File '%s' is copied from commit '%s' of source branch '%s' to target branch '%s'.",
				filePath,
				commitHash,
				sourceBranchName,
				targetBranchName,
			)
		}
	}

	return targetFile, nil
}

func (g *GitlabBranch) CreateFromDefaultBranch(verbose bool) (err error) {
	branches, err := g.GetBranches()
	if err != nil {
		return err
	}

	branchName, err := g.GetName()
	if err != nil {
		return err
	}

	_, err = branches.CreateBranchFromDefaultBranch(branchName, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabBranch) CreateMergeRequest(options *GitlabCreateMergeRequestOptions) (mergeRequest *GitlabMergeRequest, err error) {
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

	mergeRequest, err = mergeRequests.CreateMergeRequest(optionsToUse)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}

func (g *GitlabBranch) Delete(options *GitlabDeleteBranchOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
	if err != nil {
		return err
	}

	branchName, err := g.GetName()
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		return err
	}

	exists, err := g.Exists()
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

				const verboseBranchListing = false
				branchNames, err := gitlabBranches.GetBranchNames(verboseBranchListing)
				if err != nil {
					return err
				}

				exists = aslices.ContainsString(branchNames, branchName)
				if exists {
					time.Sleep(1000 * time.Millisecond)
					if options.Verbose {
						logging.LogInfof("Wait for branch '%s' to be deleted in %s .", branchName, projectUrl)
					}
				} else {
					break
				}
			}
		}

		if exists {
			return tracederrors.TracedErrorf("Internal error: failed to delete '%s' in %s", branchName, projectUrl)
		}

		if options.Verbose {
			logging.LogChangedf("Deleted branch '%s' in gitlab project %s .", branchName, projectUrl)
		}
	} else {
		if options.Verbose {
			logging.LogInfof("Branch '%s' is already absent on %s .", branchName, projectUrl)
		}
	}

	return nil
}

func (g *GitlabBranch) DeleteRepositoryFile(filePath string, commitMessage string, verbose bool) (err error) {
	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	if commitMessage == "" {
		return tracederrors.TracedErrorEmptyString("commitMessage")
	}

	fileToDelete, err := g.GetRepositoryFile(filePath, verbose)
	if err != nil {
		return err
	}

	err = fileToDelete.Delete(commitMessage, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabBranch) Exists() (exists bool, err error) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
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

func (g *GitlabBranch) GetLatestCommit(verbose bool) (latestCommit *GitlabCommit, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branchName, err := g.GetName()
	if err != nil {
		return nil, err
	}

	latestCommit, err = gitlabProject.GetLatestCommit(branchName, verbose)
	if err != nil {
		return nil, err
	}

	return latestCommit, err
}

func (g *GitlabBranch) GetLatestCommitHashAsString(verbose bool) (commitHash string, err error) {
	rawResponse, err := g.GetRawResponse()
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

	if verbose {
		logging.LogInfof("Latest commit of branch '%s' is '%s'", branchName, commitHash)
	}

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

func (g *GitlabBranch) GetNativeBranchesClientAndId() (nativeClient *gitlab.BranchesService, projectId int, err error) {
	nativeClient, err = g.GetNativeBranchesClient()
	if err != nil {
		return nil, -1, err
	}

	projectId, err = g.GetProjectId()
	if err != nil {
		return nil, -1, err
	}

	return nativeClient, projectId, nil
}

func (g *GitlabBranch) GetProjectId() (projectId int, err error) {
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

func (g *GitlabBranch) GetProjectUrl() (projectUrl string, err error) {
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

func (g *GitlabBranch) GetRawResponse() (rawResponse *gitlab.Branch, err error) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
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

func (g *GitlabBranch) GetRepositoryFile(filePath string, verbose bool) (repositoryFile *GitlabRepositoryFile, err error) {
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
		&GitlabGetRepositoryFileOptions{
			BranchName: branchName,
			Verbose:    verbose,
			Path:       filePath,
		},
	)
	if err != nil {
		return nil, err
	}

	return repositoryFile, nil
}

func (g *GitlabBranch) GetRepositoryFileSha256Sum(filePath string, verbose bool) (sha256sum string, err error) {
	if filePath == "" {
		return "", tracederrors.TracedErrorEmptyString("filePath")
	}

	repostioryFile, err := g.GetRepositoryFile(filePath, verbose)
	if err != nil {
		return "", err
	}

	sha256sum, err = repostioryFile.GetSha256CheckSum()
	if err != nil {
		return "", err
	}

	return sha256sum, nil
}

func (g *GitlabBranch) MustCopyFileToBranch(filePath string, targetBranch *GitlabBranch, verbose bool) (targetFile *GitlabRepositoryFile) {
	targetFile, err := g.CopyFileToBranch(filePath, targetBranch, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return targetFile
}

func (g *GitlabBranch) MustCreateFromDefaultBranch(verbose bool) {
	err := g.CreateFromDefaultBranch(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustCreateMergeRequest(options *GitlabCreateMergeRequestOptions) (mergeRequest *GitlabMergeRequest) {
	mergeRequest, err := g.CreateMergeRequest(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeRequest
}

func (g *GitlabBranch) MustDelete(options *GitlabDeleteBranchOptions) {
	err := g.Delete(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustDeleteRepositoryFile(filePath string, commitMessage string, verbose bool) {
	err := g.DeleteRepositoryFile(filePath, commitMessage, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustExists() (exists bool) {
	exists, err := g.Exists()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabBranch) MustGetBranches() (branches *GitlabBranches) {
	branches, err := g.GetBranches()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branches
}

func (g *GitlabBranch) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabBranch) MustGetGitlabBranches() (branches *GitlabBranches) {
	branches, err := g.GetGitlabBranches()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branches
}

func (g *GitlabBranch) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabBranch) MustGetLatestCommit(verbose bool) (latestCommit *GitlabCommit) {
	latestCommit, err := g.GetLatestCommit(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestCommit
}

func (g *GitlabBranch) MustGetLatestCommitHashAsString(verbose bool) (commitHash string) {
	commitHash, err := g.GetLatestCommitHashAsString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitHash
}

func (g *GitlabBranch) MustGetMergeRequests() (mergeRequests *GitlabProjectMergeRequests) {
	mergeRequests, err := g.GetMergeRequests()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mergeRequests
}

func (g *GitlabBranch) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabBranch) MustGetNativeBranchesClient() (nativeClient *gitlab.BranchesService) {
	nativeClient, err := g.GetNativeBranchesClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabBranch) MustGetNativeBranchesClientAndId() (nativeClient *gitlab.BranchesService, projectId int) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient, projectId
}

func (g *GitlabBranch) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabBranch) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabBranch) MustGetRawResponse() (rawResponse *gitlab.Branch) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rawResponse
}

func (g *GitlabBranch) MustGetRepositoryFile(filePath string, verbose bool) (repositoryFile *GitlabRepositoryFile) {
	repositoryFile, err := g.GetRepositoryFile(filePath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repositoryFile
}

func (g *GitlabBranch) MustGetRepositoryFileSha256Sum(filePath string, verbose bool) (sha256sum string) {
	sha256sum, err := g.GetRepositoryFileSha256Sum(filePath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sha256sum
}

func (g *GitlabBranch) MustReadFileContentAsString(options *GitlabReadFileOptions) (content string) {
	content, err := g.ReadFileContentAsString(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (g *GitlabBranch) MustRepositoryFileExists(filePath string, verbose bool) (exists bool) {
	exists, err := g.RepositoryFileExists(filePath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabBranch) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustSyncFilesToBranch(options *GitlabSyncBranchOptions) {
	err := g.SyncFilesToBranch(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustSyncFilesToBranchUsingMergeRequest(options *GitlabSyncBranchOptions) (createdMergeRequest *GitlabMergeRequest) {
	createdMergeRequest, err := g.SyncFilesToBranchUsingMergeRequest(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdMergeRequest
}

func (g *GitlabBranch) MustWriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile) {
	gitlabRepositoryFile, err := g.WriteFileContent(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabRepositoryFile
}

func (g *GitlabBranch) ReadFileContentAsString(options *GitlabReadFileOptions) (content string, err error) {
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

	content, err = gitlabProject.ReadFileContentAsString(optionsToUse)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (g *GitlabBranch) RepositoryFileExists(filePath string, verbose bool) (exists bool, err error) {
	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	repositoryFile, err := g.GetRepositoryFile(filePath, verbose)
	if err != nil {
		return false, err
	}

	exists, err = repositoryFile.Exists()
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

func (g *GitlabBranch) SyncFilesToBranch(options *GitlabSyncBranchOptions) (err error) {
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

	if options.Verbose {
		logging.LogInfof(
			"Sync files from source branch '%s' to target branch '%s' started.",
			sourceBranchName,
			targetBranchName,
		)
	}

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

	filesToSync, err := branches.GetFilesFromListWithDiffBetweenBranches(g, targetBranch, pathsToSync, options.Verbose)
	if err != nil {
		return err
	}

	if len(filesToSync) <= 0 {
		if options.Verbose {
			logging.LogInfof(
				"All '%d' files to sync from branch '%s' to target branch '%s' are already up to date.",
				len(pathsToSync),
				sourceBranchName,
				targetBranchName,
			)
		}
	} else {
		for _, pathToSync := range filesToSync {
			_, err = g.CopyFileToBranch(pathToSync, targetBranch, options.Verbose)
			if err != nil {
				return err
			}
		}

		if options.Verbose {
			logging.LogChangedf(
				"Synced '%d' files to sync from branch '%s' to target branch '%s'. '%d' files were already up to date so there was no need for a sync.",
				len(filesToSync),
				sourceBranchName,
				targetBranchName,
				len(pathsToSync)-len(filesToSync),
			)
		}
	}

	if options.Verbose {
		logging.LogInfof(
			"Sync files from source branch '%s' to target branch '%s' finished.",
			sourceBranchName,
			targetBranchName,
		)
	}

	return nil
}

func (g *GitlabBranch) SyncFilesToBranchUsingMergeRequest(options *GitlabSyncBranchOptions) (createdMergeRequest *GitlabMergeRequest, err error) {
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

	filesToSync, err := branches.GetFilesFromListWithDiffBetweenBranches(g, targetBranch, pathsToSync, options.Verbose)
	if err != nil {
		return nil, err
	}

	if len(filesToSync) > 0 {
		syncBranchName := fmt.Sprintf(
			"sync-files-%s-to-%s",
			sourceBranchName,
			targetBranchName,
		)

		optionsToUse := options.GetDeepCopy()

		err = optionsToUse.SetTargetBranchNameAndUnsetTargetBranchObject(syncBranchName)
		if err != nil {
			return nil, err
		}

		syncBranch, err := branches.CreateBranch(
			&GitlabCreateBranchOptions{
				SourceBranchName:    targetBranchName,
				BranchName:          syncBranchName,
				Verbose:             options.Verbose,
				FailIfAlreadyExists: true,
			},
		)
		if err != nil {
			return nil, err
		}

		err = g.SyncFilesToBranch(optionsToUse)
		if err != nil {
			return nil, err
		}

		courceCommitHash, err := g.GetLatestCommitHashAsString(options.Verbose)
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
			&GitlabCreateMergeRequestOptions{
				TargetBranchName:                targetBranchName,
				Title:                           mrTitle,
				Description:                     mrDescription,
				Verbose:                         options.Verbose,
				SquashEnabled:                   true,
				DeleteSourceBranchOnMerge:       true,
				FailIfMergeRequestAlreadyExists: true,
				AssignToSelf:                    true,
			},
		)
		if err != nil {
			return nil, err
		}

		if options.Verbose {
			mrUrl, err := createdMergeRequest.GetUrlAsString()
			if err != nil {
				return nil, err
			}

			logging.LogChangedf(
				"Created merge request %s to sync '%d' files from branch '%s' to target branch '%s'.",
				mrUrl,
				len(filesToSync),
				sourceBranchName,
				targetBranchName,
			)
		} else {
			logging.LogInfof(
				"All '%d' files are in sync between branch '%s' and '%s'. No merge request was created.",
				len(pathsToSync),
				syncBranchName,
				targetBranchName,
			)
		}
	}

	return createdMergeRequest, nil
}

func (g *GitlabBranch) WriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
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

	gitlabRepositoryFile, err = gitlabProject.WriteFileContent(optionsToUse)
	if err != nil {
		return nil, err
	}

	return gitlabRepositoryFile, nil
}
