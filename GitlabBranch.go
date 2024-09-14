package asciichgolangpublic

import (
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

type GitlabBranch struct {
	gitlabProject *GitlabProject
	name          string
}

func NewGitlabBranch() (g *GitlabBranch) {
	return new(GitlabBranch)
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
		return nil, TracedErrorNil("options")
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
		return TracedErrorNil("options")
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
			return TracedErrorf(
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

				exists = Slices().ContainsString(branchNames, branchName)
				if exists {
					time.Sleep(1000 * time.Millisecond)
					if options.Verbose {
						LogInfof("Wait for branch '%s' to be deleted in %s .", branchName, projectUrl)
					}
				} else {
					break
				}
			}
		}

		if exists {
			return TracedErrorf("Internal error: failed to delete '%s' in %s", branchName, projectUrl)
		}

		if options.Verbose {
			LogChangedf("Deleted branch '%s' in gitlab project %s .", branchName, projectUrl)
		}
	} else {
		if options.Verbose {
			LogInfof("Branch '%s' is already absent on %s .", branchName, projectUrl)
		}
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

		return false, TracedErrorf("Failed to evaluate if branch exists: '%w'", err)
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
		return nil, TracedErrorf("gitlabProject not set")
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
		return "", TracedError("rawCommit is nil in get latest commit hash as string")
	}

	commitHash = rawCommit.ID
	if commitHash == "" {
		return "", TracedError("commitHash is empty string after evaluation")
	}

	branchName, err := g.GetName()
	if err != nil {
		return "", err
	}

	if verbose {
		LogInfof("Latest commit of branch '%s' is '%s'", branchName, commitHash)
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
		return "", TracedErrorf("name not set")
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
		return nil, TracedErrorf("Unable to get branch: '%w'", err)
	}

	if rawResponse == nil {
		return nil, TracedError("rawResponse for GitlabBranch is nil after evaluation")
	}

	return rawResponse, nil
}

func (g *GitlabBranch) MustCreateFromDefaultBranch(verbose bool) {
	err := g.CreateFromDefaultBranch(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustCreateMergeRequest(options *GitlabCreateMergeRequestOptions) (mergeRequest *GitlabMergeRequest) {
	mergeRequest, err := g.CreateMergeRequest(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return mergeRequest
}

func (g *GitlabBranch) MustDelete(options *GitlabDeleteBranchOptions) {
	err := g.Delete(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustExists() (exists bool) {
	exists, err := g.Exists()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabBranch) MustGetBranches() (branches *GitlabBranches) {
	branches, err := g.GetBranches()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branches
}

func (g *GitlabBranch) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabBranch) MustGetGitlabBranches() (branches *GitlabBranches) {
	branches, err := g.GetGitlabBranches()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branches
}

func (g *GitlabBranch) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabBranch) MustGetLatestCommit(verbose bool) (latestCommit *GitlabCommit) {
	latestCommit, err := g.GetLatestCommit(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return latestCommit
}

func (g *GitlabBranch) MustGetLatestCommitHashAsString(verbose bool) (commitHash string) {
	commitHash, err := g.GetLatestCommitHashAsString(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitHash
}

func (g *GitlabBranch) MustGetMergeRequests() (mergeRequests *GitlabProjectMergeRequests) {
	mergeRequests, err := g.GetMergeRequests()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return mergeRequests
}

func (g *GitlabBranch) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabBranch) MustGetNativeBranchesClient() (nativeClient *gitlab.BranchesService) {
	nativeClient, err := g.GetNativeBranchesClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabBranch) MustGetNativeBranchesClientAndId() (nativeClient *gitlab.BranchesService, projectId int) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient, projectId
}

func (g *GitlabBranch) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabBranch) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabBranch) MustGetRawResponse() (rawResponse *gitlab.Branch) {
	rawResponse, err := g.GetRawResponse()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rawResponse
}

func (g *GitlabBranch) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustWriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile) {
	gitlabRepositoryFile, err := g.WriteFileContent(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabRepositoryFile
}

func (g *GitlabBranch) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabBranch) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}

func (g *GitlabBranch) WriteFileContent(options *GitlabWriteFileOptions) (gitlabRepositoryFile *GitlabRepositoryFile, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
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
