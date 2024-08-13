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

func (g *GitlabBranch) Delete(verbose bool) (err error) {
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

		// Deleting is not instantaneous so lets check if deleted branch is really absent:
		for i := 0; i < 10; i++ {
			exists, err = g.Exists()
			if err != nil {
				return err
			}

			if exists {
				time.Sleep(500 * time.Millisecond)
				if verbose {
					LogInfof("Wait for branch '%s' to be deleted in %s .", branchName, projectUrl)
				}
			}
		}

		if exists {
			return TracedErrorf("Internal error: failed to dete '%s' in %s", branchName, projectUrl)
		}

		if verbose {
			LogChangedf("Deleted branch '%s' in gitlab project %s .", branchName, projectUrl)
		}
	} else {
		if verbose {
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

func (g *GitlabBranch) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
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

func (g *GitlabBranch) MustCreateFromDefaultBranch(verbose bool) {
	err := g.CreateFromDefaultBranch(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabBranch) MustDelete(verbose bool) {
	err := g.Delete(verbose)
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

func (g *GitlabBranch) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
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
