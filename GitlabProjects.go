package asciichgolangpublic

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var ErrGitlabProjectNotFound = errors.New("Gitlab project not found")

type GitlabProjects struct {
	gitlab *GitlabInstance
}

func NewGitlabProjects() (gitlabProject *GitlabProjects) {
	return new(GitlabProjects)
}

func (g *GitlabProject) DeleteAllRepositoryFiles(ctx context.Context, branchName string) (err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return err
	}

	err = repositoryFiles.DeleteAllRepositoryFiles(ctx, branchName)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabProject) HasNoRepositoryFiles(ctx context.Context, branchName string) (hasNoRepositoryFiles bool, err error) {
	repositoryFiles, err := g.GetRepositoryFiles()
	if err != nil {
		return false, err
	}

	hasNoRepositoryFiles, err = repositoryFiles.HasNoRepositoryFiles(ctx, branchName)
	if err != nil {
		return false, err
	}

	return hasNoRepositoryFiles, nil
}

func (g *GitlabProjects) DeleteProject(ctx context.Context, deleteProjectOptions *GitlabDeleteProjectOptions) (err error) {
	if deleteProjectOptions == nil {
		return tracederrors.TracedErrorNil("deleteProjectOptions")
	}

	fqdn, err := g.GetFqdn()
	if err != nil {
		return err
	}

	projectPath, err := deleteProjectOptions.GetProjectPath()
	if err != nil {
		return err
	}

	projectExists, err := g.ProjectByProjectPathExists(ctx, projectPath)
	if err != nil {
		return err
	}

	if projectExists {
		nativeProjectsService, err := g.GetNativeProjectsService()
		if err != nil {
			return err
		}

		isPersonalProject, err := g.IsProjectPathPersonalProject(projectPath)
		if err != nil {
			return err
		}

		if isPersonalProject {
			personalProjectsPath, err := g.GetPersonalProjectsPath(ctx)
			if err != nil {
				return err
			}

			currentUserName, err := g.GetCurrentUserName(ctx)
			if err != nil {
				return err
			}

			projectPath = fmt.Sprintf(
				"%s/%s",
				currentUserName,
				strings.TrimPrefix(projectPath, personalProjectsPath),
			)
			projectPath = strings.ReplaceAll(projectPath, "//", "/")
		}

		_, err = nativeProjectsService.DeleteProject(projectPath, &gitlab.DeleteProjectOptions{})
		if err != nil {
			return tracederrors.TracedErrorf(
				"Failed to delete gitlab project '%s' on instance '%s': '%w'",
				projectPath,
				fqdn,
				err,
			)
		}

		logging.LogChangedByCtxf(ctx, "Delete project '%s' on gitlab '%s'.", projectPath, fqdn)
	} else {
		logging.LogInfoByCtxf(ctx, "Project '%s' is already absent on gitlab '%s'.", projectPath, fqdn)
	}

	return nil
}

func (g *GitlabProjects) GetCurrentUserName(ctx context.Context) (userName string, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return "", err
	}

	userName, err = gitlab.GetCurrentUsersName(ctx)
	if err != nil {
		return "", err
	}

	return userName, nil
}

func (g *GitlabProjects) GetPersonalProjectsPath(ctx context.Context) (personalProjectPath string, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return "", err
	}

	personalProjectPath, err = gitlab.GetPersonalProjectsPath(ctx)
	if err != nil {
		return "", err
	}

	return personalProjectPath, nil
}

func (g *GitlabProjects) GetProjectById(projectId int) (gitlabProject *GitlabProject, err error) {
	if projectId <= 0 {
		return nil, tracederrors.TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
	}

	nativeProjectsClient, err := g.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	nativeProject, _, err := nativeProjectsClient.GetProject(projectId, &gitlab.GetProjectOptions{})
	if err != nil {
		if stringsutils.ContainsAtLeastOneSubstring(err.Error(), []string{"404 {message: 404 Project Not Found}"}) {
			return nil, tracederrors.TracedErrorf("%w: %d", ErrGitlabProjectNotFound, projectId)
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
		return nil, tracederrors.TracedErrorNil("nativeProject")
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

func (g *GitlabProjects) GetProjectByProjectPath(ctx context.Context, projectPath string) (gitlabProject *GitlabProject, err error) {
	if len(projectPath) <= 0 {
		return nil, tracederrors.TracedError("projectPath is empty string")
	}

	nativeProjectsClient, err := g.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	isPersonalProject, err := g.IsProjectPathPersonalProject(projectPath)
	if err != nil {
		return nil, err
	}

	if isPersonalProject {
		ownedProjects, err := g.GetProjectList(
			&GitlabgetProjectListOptions{
				Verbose: false,
				Owned:   true,
			},
		)
		if err != nil {
			return nil, err
		}

		currentUserName, err := g.GetCurrentUserName(ctx)
		if err != nil {
			return nil, err
		}

		personalProjectsPath, err := g.GetPersonalProjectsPath(ctx)
		if err != nil {
			return nil, err
		}

		expectedPrivatePath := fmt.Sprintf(
			"%s/%s",
			currentUserName,
			strings.TrimPrefix(projectPath, personalProjectsPath),
		)
		expectedPrivatePath = strings.ReplaceAll(expectedPrivatePath, "//", "/")

		for _, toCheck := range ownedProjects {
			pathToCheck, err := toCheck.GetCachedPath(ctx)
			if err != nil {
				return nil, err
			}

			if expectedPrivatePath == pathToCheck {
				gitlabProject = toCheck
				break
			}
		}

		if gitlabProject == nil {
			errorNotFound := tracederrors.TracedErrorf("%w: Personal project %s", ErrGitlabProjectNotFound, projectPath)
			return nil, errorNotFound
		}
	} else {
		nativeProject, _, err := nativeProjectsClient.GetProject(projectPath, &gitlab.GetProjectOptions{})
		if err != nil {

			if stringsutils.ContainsAtLeastOneSubstring(err.Error(), []string{"404 {message: 404 Project Not Found}", "404 Not Found"}) {
				errorNotFound := tracederrors.TracedErrorf("%w: %s", ErrGitlabProjectNotFound, projectPath)
				return nil, errorNotFound
			}
			return nil, err
		}

		gitlabProject, err = g.GetProjectByNativeProject(nativeProject)
		if err != nil {
			return nil, err
		}
	}

	return gitlabProject, nil
}

func (g *GitlabProjects) GetProjectIdByProjectPath(ctx context.Context, projectPath string) (projectId int, err error) {
	if projectPath == "" {
		return -1, tracederrors.TracedErrorEmptyString("projectPath")
	}

	project, err := g.GetProjectByProjectPath(ctx, projectPath)
	if err != nil {
		return -1, err
	}

	projectId, err = project.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabProjects) ProjectByProjectIdExists(ctx context.Context, projectId int) (projectExists bool, err error) {
	if projectId <= 0 {
		return false, tracederrors.TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
	}

	_, err = g.GetProjectById(projectId)
	if err != nil {
		if errors.Is(err, ErrGitlabGroupNotFoundError) {
			logging.LogInfoByCtxf(ctx, "Gitlab project with id '%d' does not exist.", projectId)
			return false, nil
		}
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Gitlab project with id '%d' does exist.", projectId)
	return true, nil
}

func (p *GitlabProjects) CreateProject(ctx context.Context, createOptions *GitlabCreateProjectOptions) (createdGitlabProject *GitlabProject, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	fqdn, err := p.GetFqdn()
	if err != nil {
		return nil, err
	}

	projectPath, err := createOptions.GetProjectPath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create project '%s' on gitlab '%s' started.", projectPath, fqdn)

	projectExists, err := p.ProjectByProjectPathExists(ctx, projectPath)
	if err != nil {
		return nil, err
	}

	if projectExists {
		logging.LogInfoByCtxf(ctx, "Gitlab project '%s' on gitlab '%s' already exists.", projectPath, fqdn)
	} else {
		logging.LogInfoByCtxf(ctx, "Going to create gitlab project '%s' on '%s'.", projectPath, fqdn)

		isPersonalProject, err := p.IsProjectPathPersonalProject(projectPath)
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
		if !isPersonalProject {
			groupPath, err := createOptions.GetGroupPath(ctx)
			if err != nil {
				return nil, err
			}

			logging.LogInfoByCtxf(ctx, "groupPath for creating gitlab project '%s' is '%s'.", projectPath, groupPath)

			asciichgolangGitlab, err := p.GetGitlab()
			if err != nil {
				return nil, err
			}

			if groupPath != "" {
				createdGroup, err := asciichgolangGitlab.CreateGroupByPath(
					ctx,
					groupPath,
				)
				if err != nil {
					return nil, err
				}

				groupIdForNewProject, err = createdGroup.GetId(ctx)
				if err != nil {
					return nil, err
				}

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

		logging.LogChangedByCtxf(ctx, "Created gitlab project '%s' on '%s'.", projectPath, fqdn)
	}

	createdGitlabProject, err = p.GetProjectByProjectPath(ctx, projectPath)
	if err != nil {
		return nil, err
	}

	if createOptions.IsPublic {
		err = createdGitlabProject.MakePublic(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		err = createdGitlabProject.MakePrivate(ctx)
		if err != nil {
			return nil, err
		}
	}

	logging.LogInfoByCtxf(ctx, "Create project '%s' on gitlab '%s' finished.", projectPath, fqdn)

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
		return nil, tracederrors.TracedError("gitlab is not set")
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
		return nil, tracederrors.TracedError("Unable to get nativeGitlabProject")
	}

	return nativeGitlabProject, nil
}

func (p *GitlabProjects) GetProjectList(options *GitlabgetProjectListOptions) (gitlabProjects []*GitlabProject, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

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
			return nil, tracederrors.TracedErrorf("Unable  to get gitlab native project list: '%w'", err)
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

func (p *GitlabProjects) GetProjectPathList(ctx context.Context, options *GitlabgetProjectListOptions) (projectPaths []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	projects, err := p.GetProjectList(options)
	if err != nil {
		return nil, err
	}

	projectPaths = []string{}
	for _, nativeProject := range projects {
		pathToAdd, err := nativeProject.GetCachedPath(ctx)
		if err != nil {
			return nil, err
		}

		projectPaths = append(projectPaths, pathToAdd)
	}

	sort.Strings(projectPaths)

	return projectPaths, nil
}

func (p *GitlabProjects) IsProjectPathPersonalProject(projectPath string) (isPersonalProject bool, err error) {
	if projectPath == "" {
		return false, tracederrors.TracedErrorEmptyString("projectPath")
	}

	isPersonalProject = stringsutils.HasAtLeastOnePrefix(projectPath, []string{"users/", "/users/"})

	return isPersonalProject, nil
}

func (p *GitlabProjects) ProjectByProjectPathExists(ctx context.Context, projectPath string) (projectExists bool, err error) {
	if len(projectPath) <= 0 {
		return false, tracederrors.TracedError("projectPath is empty string")
	}

	_, err = p.GetProjectByProjectPath(ctx, projectPath)
	if err != nil {
		if errors.Is(err, ErrGitlabProjectNotFound) {
			logging.LogInfoByCtxf(ctx, "Gitlab project '%s' does not exist.", projectPath)
			return false, nil
		}
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Gitlab project '%s' does exist.", projectPath)
	return true, nil
}

func (p *GitlabProjects) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}
