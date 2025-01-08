package asciichgolangpublic

import (
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

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

func (g *GitlabBranches) CreateBranch(options *GitlabCreateBranchOptions) (createdBranch *GitlabBranch, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	nativeClient, projectId, err := g.GetNativeBranchesClientAndId()
	if err != nil {
		return nil, err
	}

	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		return nil, err
	}

	branchName, err := options.GetBranchName()
	if err != nil {
		return nil, err
	}

	exists, err := g.BranchByNameExists(branchName)
	if err != nil {
		return nil, err
	}

	sourceBranchName, err := options.GetSourceBranchName()
	if err != nil {
		return nil, err
	}

	if exists {
		if options.FailIfAlreadyExists {
			return nil, TracedErrorf(
				"Branch '%s' already exists in gitlab project %s .", branchName, projectUrl,
			)
		}

		if options.Verbose {
			LogInfof("Branch '%s' already exists in gitlab project %s .", branchName, projectUrl)
		}
	} else {
		_, _, err = nativeClient.CreateBranch(
			projectId,
			&gitlab.CreateBranchOptions{
				Branch: &branchName,
				Ref:    &sourceBranchName,
			},
		)
		if err != nil {
			return nil, TracedErrorf(
				"Unable to create branch '%s' from branch '%s' in gitlab project %s : '%w'",
				branchName,
				sourceBranchName,
				projectUrl,
				err,
			)
		}

		if options.Verbose {
			LogChangedf(
				"Created branch '%s' from '%s' in gitlab project %s .",
				branchName,
				sourceBranchName,
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

	createdBranch, err = g.CreateBranch(
		&GitlabCreateBranchOptions{
			SourceBranchName: sourceBranch,
			BranchName:       branchName,
			Verbose:          verbose,
		},
	)
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

	deletedBranchNames := []string{}
	for _, toDelete := range branches {
		err = toDelete.Delete(&GitlabDeleteBranchOptions{
			SkipWaitForDeletion: true,
			Verbose:             verbose,
		})
		if err != nil {
			return err
		}

		branchName, err := toDelete.GetName()
		if err != nil {
			return err
		}

		deletedBranchNames = append(deletedBranchNames, branchName)
	}

	branchNotDeletedYetFound := false
	for i := 0; i < 30; i++ {
		branchNotDeletedYetFound = false

		const verboseList bool = false
		currentBranchNames, err := g.GetBranchNames(verboseList)
		if err != nil {
			return err
		}

		for _, deleted := range deletedBranchNames {
			if Slices().ContainsString(currentBranchNames, deleted) {
				branchNotDeletedYetFound = true
				break
			}
		}

		if branchNotDeletedYetFound {
			if verbose {
				LogInfof("Wait for all non default branches to be deleted.")
				time.Sleep(1 * time.Second)
			}
		} else {
			break
		}
	}

	if branchNotDeletedYetFound {
		return TracedError("Unable to delete all branches except default branch")
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

func (g *GitlabBranches) GetFilesFromListWithDiffBetweenBranches(branchA *GitlabBranch, branchB *GitlabBranch, filesToCheck []string, verbose bool) (filesWithDiffBetweenBranches []string, err error) {
	if branchA == nil {
		return nil, TracedErrorNil("branchA")
	}

	if branchB == nil {
		return nil, TracedErrorNil("branchB")
	}

	if len(filesToCheck) <= 0 {
		return nil, TracedError("filesToCkeck has no elements")
	}

	branchAName, err := branchA.GetName()
	if err != nil {
		return nil, err
	}

	branchBName, err := branchB.GetName()
	if err != nil {
		return nil, err
	}

	for _, toCheck := range filesToCheck {
		checksumA, err := branchA.GetRepositoryFileSha256Sum(toCheck, verbose)
		if err != nil {
			return nil, err
		}

		checksumB := ""

		targetFileExists, err := branchB.RepositoryFileExists(toCheck, verbose)
		if err != nil {
			return nil, err
		}

		if targetFileExists {
			checksumB, err = branchB.GetRepositoryFileSha256Sum(toCheck, verbose)
			if err != nil {
				return nil, err
			}
		}

		if checksumA == checksumB {
			if verbose {
				LogInfof(
					"File '%s' in branch '%s' and '%s' is equal with sha256sum '%s'.",
					toCheck,
					branchAName,
					branchBName,
					checksumA,
				)
			}
			continue
		}

		if verbose {
			if targetFileExists {
				LogInfof(
					"File '%s' in branch '%s' has sha256sum '%s' and does not exist in branchB '%s'. This is considered a difference.",
					toCheck,
					branchAName,
					checksumA,
					branchBName,
				)
			} else {
				LogInfof(
					"File '%s in branch '%s' has sha256sum '%s' and is not equal to branch '%s' where sha256sum is '%s'.",
					toCheck,
					branchAName,
					checksumA,
					branchBName,
					checksumB,
				)
			}
		}

		filesWithDiffBetweenBranches = append(filesWithDiffBetweenBranches, toCheck)
	}

	if verbose {
		if len(filesWithDiffBetweenBranches) > 0 {
			LogInfof(
				"Found '%d' out of '%d' files with different content between branch '%s' and '%s'.",
				len(filesWithDiffBetweenBranches),
				len(filesToCheck),
				branchAName,
				branchBName,
			)
		} else {
			LogInfof(
				"All '%d' files of branch '%s' and '%s' have equal content.",
				len(filesToCheck),
				branchAName,
				branchBName,
			)
		}
	}

	return filesWithDiffBetweenBranches, nil
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

func (g *GitlabBranches) MustCreateBranch(options *GitlabCreateBranchOptions) (createdBranch *GitlabBranch) {
	createdBranch, err := g.CreateBranch(options)
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

func (g *GitlabBranches) MustGetFilesFromListWithDiffBetweenBranches(branchA *GitlabBranch, branchB *GitlabBranch, filesToCheck []string, verbose bool) (filesWithDiffBetweenBranches []string) {
	filesWithDiffBetweenBranches, err := g.GetFilesFromListWithDiffBetweenBranches(branchA, branchB, filesToCheck, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filesWithDiffBetweenBranches
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
