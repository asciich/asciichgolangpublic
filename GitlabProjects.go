package asciichgolangpublic

import (
	"errors"

	"github.com/xanzy/go-gitlab"
)

var ErrGitlabProjectNotFound = errors.New("Gitlab project not found")

type GitlabProjects struct {
	gitlab *GitlabInstance
}

func NewGitlabProjects() (gitlabProject *GitlabProjects) {
	return new(GitlabProjects)
}

func (g *GitlabProjects) GetProjectById(projectId int) (gitlabProject *GitlabProject, err error) {
	if projectId <= 0 {
		return nil, TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
	}

	nativeProjectsClient, err := g.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	nativeProject, _, err := nativeProjectsClient.GetProject(projectId, &gitlab.GetProjectOptions{})
	if err != nil {
		if Strings().ContainsAtLeastOneSubstring(err.Error(), []string{"404 {message: 404 Project Not Found}"}) {
			return nil, TracedErrorf("%w: %d", ErrGitlabProjectNotFound, projectId)
		}
		return nil, err
	}

	gitlabProject, err = g.GetProjectByNativeProject(nativeProject)
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabProjects) GetProjectByNativeProject(nativeProject *gitlab.Project) (gitlabProject *GitlabProject, err error) {
	if nativeProject == nil {
		return nil, TracedErrorNil("nativeProject")
	}

	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	gitlabProject = NewGitlabProject()
	err = gitlabProject.SetGitlab(gitlab)
	if err != nil {
		return nil, err
	}

	err = gitlabProject.SetId(nativeProject.ID)
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabProjects) GetProjectByProjectPath(projectPath string) (gitlabProject *GitlabProject, err error) {
	if len(projectPath) <= 0 {
		return nil, TracedError("projectPath is empty string")
	}

	nativeProjectsClient, err := g.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	nativeProject, _, err := nativeProjectsClient.GetProject(projectPath, &gitlab.GetProjectOptions{})
	if err != nil {
		if Strings().ContainsAtLeastOneSubstring(err.Error(), []string{"404 {message: 404 Project Not Found}", "404 Not Found"}) {
			return nil, TracedErrorf("%w: %s", ErrGitlabProjectNotFound, projectPath)
		}
		return nil, err
	}

	gitlabProject, err = g.GetProjectByNativeProject(nativeProject)
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabProjects) MustCreateProject(createOptions *GitlabCreateProjectOptions) (createdGitlabProject *GitlabProject) {
	createdGitlabProject, err := g.CreateProject(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdGitlabProject
}

func (g *GitlabProjects) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabProjects) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabProjects) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabProjects) MustGetNativeProjectsService() (nativeGitlabProject *gitlab.ProjectsService) {
	nativeGitlabProject, err := g.GetNativeProjectsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeGitlabProject
}

func (g *GitlabProjects) MustGetProjectById(projectId int) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetProjectById(projectId)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabProjects) MustGetProjectByNativeProject(nativeProject *gitlab.Project) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetProjectByNativeProject(nativeProject)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabProjects) MustGetProjectByProjectPath(projectPath string) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetProjectByProjectPath(projectPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabProjects) MustGetProjectList(verbose bool) (gitlabProjects []*GitlabProject) {
	gitlabProjects, err := g.GetProjectList(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProjects
}

func (g *GitlabProjects) MustGetProjectPathList(verbose bool) (projectPaths []string) {
	projectPaths, err := g.GetProjectPathList(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectPaths
}

func (g *GitlabProjects) MustProjectByProjectIdExists(projectId int, verbose bool) (projectExists bool) {
	projectExists, err := g.ProjectByProjectIdExists(projectId, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabProjects) MustProjectByProjectPathExists(projectPath string, verbose bool) (projectExists bool) {
	projectExists, err := g.ProjectByProjectPathExists(projectPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabProjects) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabProjects) ProjectByProjectIdExists(projectId int, verbose bool) (projectExists bool, err error) {
	if projectId <= 0 {
		return false, TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
	}

	_, err = g.GetProjectById(projectId)
	if err != nil {
		if errors.Is(err, ErrGitlabGroupNotFoundError) {
			if verbose {
				LogInfof("Gitlab project with id '%d' does not exist.", projectId)
			}
			return false, nil
		}
		return false, err
	}

	if verbose {
		LogInfof("Gitlab project with id '%d' does exist.", projectId)
	}
	return true, nil
}

func (p *GitlabProjects) CreateProject(createOptions *GitlabCreateProjectOptions) (createdGitlabProject *GitlabProject, err error) {
	if createOptions == nil {
		return nil, TracedError("createOptions is nil")
	}

	fqdn, err := p.GetFqdn()
	if err != nil {
		return nil, err
	}

	projectPath, err := createOptions.GetProjectPath()
	if err != nil {
		return nil, err
	}

	if createOptions.Verbose {
		LogInfof("Create project '%s' on gitlab '%s' started.", projectPath, fqdn)
	}

	projectExists, err := p.ProjectByProjectPathExists(projectPath, createOptions.Verbose)
	if err != nil {
		return nil, err
	}

	if projectExists {
		if createOptions.Verbose {
			LogInfof("Gitlab project '%s' on gitlab '%s' already exists.", projectPath, fqdn)
		}
	} else {
		if createOptions.Verbose {
			LogInfof("Going to create gitlab project '%s' on '%s'.", projectPath, fqdn)
		}

		groupPath, err := createOptions.GetGroupPath(createOptions.Verbose)
		if err != nil {
			return nil, err
		}

		if createOptions.Verbose {
			LogInfof("groupPath for creationg gitlab project '%s' is '%s'.", projectPath, groupPath)
		}

		asciichgolangGitlab, err := p.GetGitlab()
		if err != nil {
			return nil, err
		}

		nativeProjectsService, err := p.GetNativeProjectsService()
		if err != nil {
			return nil, err
		}

		projectName, err := createOptions.GetProjectName()
		if err != nil {
			return nil, err
		}

		groupIdForNewProject := -1
		if groupPath != "" {
			createdGroup, err := asciichgolangGitlab.CreateGroup(
				&GitlabCreateGroupOptions{
					GroupPath: groupPath,
					Verbose:   createOptions.Verbose,
				},
			)
			if err != nil {
				return nil, err
			}

			groupIdForNewProject, err = createdGroup.GetId()
			if err != nil {
				return nil, err
			}

		}

		createProjectOptions := &gitlab.CreateProjectOptions{
			Name: &projectName,
		}

		if groupIdForNewProject > 0 {
			createProjectOptions.NamespaceID = &groupIdForNewProject
		}

		_, _, err = nativeProjectsService.CreateProject(
			createProjectOptions,
		)
		if err != nil {
			return nil, err
		}

		if createOptions.Verbose {
			LogChangedf("Created gitlab project '%s' on '%s'.", projectPath, fqdn)
		}
	}

	createdGitlabProject, err = p.GetProjectByProjectPath(projectPath)
	if err != nil {
		return nil, err
	}

	if createOptions.IsPublic {
		err = createdGitlabProject.MakePublic(createOptions.Verbose)
		if err != nil {
			return nil, err
		}
	} else {
		err = createdGitlabProject.MakePrivate(createOptions.Verbose)
		if err != nil {
			return nil, err
		}
	}

	if createOptions.Verbose {
		LogInfof("Create project '%s' on gitlab '%s' finished.", projectPath, fqdn)
	}

	return createdGitlabProject, nil
}

func (p *GitlabProjects) GetFqdn() (fqdn string, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return "", err
	}

	fqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (p *GitlabProjects) GetGitlab() (gitlab *GitlabInstance, err error) {
	if p.gitlab == nil {
		return nil, TracedError("gitlab is not set")
	}

	return p.gitlab, nil
}

func (p *GitlabProjects) GetNativeClient() (nativeClient *gitlab.Client, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (p *GitlabProjects) GetNativeProjectsService() (nativeGitlabProject *gitlab.ProjectsService, err error) {
	nativeClient, err := p.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeGitlabProject = nativeClient.Projects
	if nativeGitlabProject == nil {
		return nil, TracedError("Unable to get nativeGitlabProject")
	}

	return nativeGitlabProject, nil
}

func (p *GitlabProjects) GetProjectList(verbose bool) (gitlabProjects []*GitlabProject, err error) {
	nativeService, err := p.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	var nativeList []*gitlab.Project
	pageNumber := 1
	for {
		partList, response, err := nativeService.ListProjects(
			&gitlab.ListProjectsOptions{
				ListOptions: gitlab.ListOptions{
					Page:    pageNumber,
					PerPage: 50,
				},
			},
		)
		if err != nil {
			return nil, TracedErrorf("Unable  to get gitlab native project list: '%w'", err)
		}

		nativeList = append(nativeList, partList...)
		if response.NextPage <= 0 {
			break
		}

		pageNumber = response.NextPage
	}

	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	gitlabProjects = []*GitlabProject{}
	for _, nativeProject := range nativeList {
		projectToAdd := NewGitlabProject()
		if err != nil {
			return nil, err
		}

		err = projectToAdd.SetGitlab(gitlab)
		if err != nil {
			return nil, err
		}

		err = projectToAdd.SetId(nativeProject.ID)
		if err != nil {
			return nil, err
		}

		err = projectToAdd.SetCachedPath(nativeProject.PathWithNamespace)
		if err != nil {
			return nil, err
		}

		gitlabProjects = append(gitlabProjects, projectToAdd)
	}

	return gitlabProjects, nil
}

func (p *GitlabProjects) GetProjectPathList(verbose bool) (projectPaths []string, err error) {
	projects, err := p.GetProjectList(verbose)
	if err != nil {
		return nil, err
	}

	projectPaths = []string{}
	for _, nativeProject := range projects {
		pathToAdd, err := nativeProject.GetCachedPath()
		if err != nil {
			return nil, err
		}

		projectPaths = append(projectPaths, pathToAdd)
	}

	projectPaths = Slices().SortStringSlice(projectPaths)

	return projectPaths, nil
}

func (p *GitlabProjects) ProjectByProjectPathExists(projectPath string, verbose bool) (projectExists bool, err error) {
	if len(projectPath) <= 0 {
		return false, TracedError("projectPath is empty string")
	}

	_, err = p.GetProjectByProjectPath(projectPath)
	if err != nil {
		if errors.Is(err, ErrGitlabProjectNotFound) {
			if verbose {
				LogInfof("Gitlab project '%s' does not exist.", projectPath)
			}
			return false, nil
		}
		return false, err
	}

	if verbose {
		LogInfof("Gitlab project '%s' does exist.", projectPath)
	}
	return true, nil
}

func (p *GitlabProjects) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
