package asciichgolangpublic

import "github.com/xanzy/go-gitlab"

type GitlabBranches struct {
	gitlabProject *GitlabProject
}

func NewGitlabBranches() (g *GitlabBranches) {
	return new(GitlabBranches)
}

func (g *GitlabBranches) BranchByNameExists(branchName string) (exists bool, err error) {
	if branchName == "" {
		return false, TracedErrorEmptyString("branchName")
	}

	branch, err := g.GetBranchByName(branchName)
	if err != nil {
		return false, err
	}

	exists, err = branch.Exists()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (g *GitlabBranches) CreateBranch(sourceBranch string, branchName string, verbose bool) (createdBranch *GitlabBranch, err error) {
	if sourceBranch == "" {
		return nil, TracedErrorEmptyString("sourceBranch")
	}

	if branchName == "" {
		return nil, TracedErrorEmptyString("branchName")
	}

	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
	if err != nil {
		return nil, err
	}

	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		return nil, err
	}

	exists, err := g.BranchByNameExists(branchName)
	if err != nil {
		return nil, err
	}

	if exists {
		if verbose {
			LogInfof("Branch '%s' already exists in gitlab project %s .", branchName, projectUrl)
		}
	} else {
		_, _, err = nativeClient.CreateBranch(
			projectId,
			&gitlab.CreateBranchOptions{
				Branch: &branchName,
				Ref:    &sourceBranch,
			},
		)
		if err != nil {
			return nil, TracedErrorf(
				"Unable to create branch '%s' from branch '%s' in gitlab project %s : '%w'",
				branchName,
				sourceBranch,
				projectUrl,
				err,
			)
		}

		if verbose {
			LogChangedf(
				"Created branch '%s' from '%s' in gitlab project %s .",
				branchName,
				sourceBranch,
				projectUrl,
			)
		}
	}

	createdBranch, err = g.GetBranchByName(branchName)
	if err != nil {
		return nil, err
	}

	return createdBranch, nil
}

func (g *GitlabBranches) CreateBranchFromDefaultBranch(branchName string, verbose bool) (createdBranch *GitlabBranch, err error) {
	if branchName == "" {
		return nil, TracedErrorEmptyString("branchName")
	}

	sourceBranch, err := g.GetDefaultBranchName()
	if err != nil {
		return nil, err
	}

	createdBranch, err = g.CreateBranch(sourceBranch, branchName, verbose)
	if err != nil {
		return nil, err
	}

	return createdBranch, nil
}

func (g *GitlabBranches) DeleteAllBranchesExceptDefaultBranch(verbose bool) (err error) {
	branches, err := g.GetBranchesExceptDefaultBranch(verbose)
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		return err
	}

	for _, toDelete := range branches {
		err = toDelete.Delete(verbose)
		if err != nil {
			return err
		}
	}

	if len(branches) > 0 {
		LogChangedf("Deleted '%d' branches from gitlab project %s .", len(branches), projectUrl)
	} else {
		LogInfof("No branches to delete in gitlab project %s .", projectUrl)
	}

	return nil
}

func (g *GitlabBranches) GetBranchByName(branchName string) (branch *GitlabBranch, err error) {
	if branchName == "" {
		return nil, TracedErrorNil("branchName")
	}

	project, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branch = NewGitlabBranch()

	err = branch.SetGitlabProject(project)
	if err != nil {
		return nil, err
	}

	err = branch.SetName(branchName)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (g *GitlabBranches) GetBranchNames(verbose bool) (branchNames []string, err error) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
	if err != nil {
		return nil, err
	}

	nextPage := 1

	branchNames = []string{}

	for {
		if nextPage <= 0 {
			break
		}

		nativeBranches, response, err := nativeClient.ListBranches(
			projectId,
			&gitlab.ListBranchesOptions{
				ListOptions: gitlab.ListOptions{
					Page: nextPage,
				},
			},
		)
		if err != nil {
			return nil, TracedErrorf("Unable to get branch list: '%w'", err)
		}

		for _, toAdd := range nativeBranches {
			branchNames = append(branchNames, toAdd.Name)
		}

		nextPage = response.NextPage
	}

	return branchNames, nil
}

func (g *GitlabBranches) GetBranchNamesExceptDefaultBranch(verbose bool) (branchNames []string, err error) {
	allBranchNames, err := g.GetBranchNames(verbose)
	if err != nil {
		return nil, err
	}

	defaultBranchName, err := g.GetDefaultBranchName()
	if err != nil {
		return nil, err
	}

	branchNames = Slices().RemoveString(allBranchNames, defaultBranchName)

	return branchNames, nil
}

func (g *GitlabBranches) GetBranchesExceptDefaultBranch(verbose bool) (branches []*GitlabBranch, err error) {
	branchNames, err := g.GetBranchNamesExceptDefaultBranch(verbose)
	if err != nil {
		return nil, err
	}

	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	branches = []*GitlabBranch{}
	for _, name := range branchNames {
		toAdd := NewGitlabBranch()

		err = toAdd.SetGitlabProject(gitlabProject)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetName(name)
		if err != nil {
			return nil, err
		}

		branches = append(branches, toAdd)
	}

	return branches, nil
}

func (g *GitlabBranches) GetDefaultBranchName() (defaultBranchName string, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	defaultBranchName, err = gitlabProject.GetDefaultBranchName()
	if err != nil {
		return "", err
	}

	return defaultBranchName, nil
}

func (g *GitlabBranches) GetGitlab() (gitlab *GitlabInstance, err error) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	gitlab, err = gitlabProject.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func (g *GitlabBranches) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabBranches) GetNativeBranchesClient() (nativeBranches *gitlab.BranchesService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeBranches, err = gitlab.GetNativeBranchesClient()
	if err != nil {
		return nil, err
	}

	return nativeBranches, nil
}

func (g *GitlabBranches) GetNativeBranchesClientAndId() (nativeClient *gitlab.BranchesService, projectId int, err error) {
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

func (g *GitlabBranches) GetProjectId() (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return 1, err
	}

	projectId, err = project.GetId()
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabBranches) GetProjectUrl() (projectUrl string, err error) {
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

func (g *GitlabBranches) MustBranchByNameExists(branchName string) (exists bool) {
	exists, err := g.BranchByNameExists(branchName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabBranches) MustCreateBranch(sourceBranch string, branchName string, verbose bool) (createdBranch *GitlabBranch) {
	createdBranch, err := g.CreateBranch(sourceBranch, branchName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdBranch
}

func (g *GitlabBranches) MustCreateBranchFromDefaultBranch(branchName string, verbose bool) (createdBranch *GitlabBranch) {
	createdBranch, err := g.CreateBranchFromDefaultBranch(branchName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdBranch
}

func (g *GitlabBranches) MustDeleteAllBranchesExceptDefaultBranch(verbose bool) {
	err := g.DeleteAllBranchesExceptDefaultBranch(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabBranches) MustGetBranchByName(branchName string) (branch *GitlabBranch) {
	branch, err := g.GetBranchByName(branchName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branch
}

func (g *GitlabBranches) MustGetBranchNames(verbose bool) (branchNames []string) {
	branchNames, err := g.GetBranchNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branchNames
}

func (g *GitlabBranches) MustGetBranchNamesExceptDefaultBranch(verbose bool) (branchNames []string) {
	branchNames, err := g.GetBranchNamesExceptDefaultBranch(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branchNames
}

func (g *GitlabBranches) MustGetBranchesExceptDefaultBranch(verbose bool) (branches []*GitlabBranch) {
	branches, err := g.GetBranchesExceptDefaultBranch(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branches
}

func (g *GitlabBranches) MustGetDefaultBranchName() (defaultBranchName string) {
	defaultBranchName, err := g.GetDefaultBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return defaultBranchName
}

func (g *GitlabBranches) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabBranches) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabBranches) MustGetNativeBranchesClient() (nativeBranches *gitlab.BranchesService) {
	nativeBranches, err := g.GetNativeBranchesClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeBranches
}

func (g *GitlabBranches) MustGetNativeBranchesClientAndId() (nativeClient *gitlab.BranchesService, projectId int) {
	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient, projectId
}

func (g *GitlabBranches) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabBranches) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabBranches) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabBranches) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}
