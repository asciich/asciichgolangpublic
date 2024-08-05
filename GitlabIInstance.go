package asciichgolangpublic

import (
	"fmt"
	"time"

	"github.com/xanzy/go-gitlab"
)

type GitlabInstance struct {
	fqdn                     string
	nativeClient             *gitlab.Client
	currentlyUsedAccessToken *string
}

func GetGitlabByFQDN(fqdn string) (gitlab *GitlabInstance, err error) {
	if len(fqdn) <= 0 {
		return nil, TracedError("fqdn is empty string")
	}

	gitlab = NewGitlab()
	err = gitlab.SetFqdn(fqdn)
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func MustGetGitlabByFQDN(fqdn string) (gitlab *GitlabInstance) {
	gitlab, err := GetGitlabByFQDN(fqdn)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func MustGetGitlabByFqdn(fqdn string) (gitlab *GitlabInstance) {
	gitlab, err := GetGitlabByFQDN(fqdn)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func NewGitlab() (gitlab *GitlabInstance) {
	return new(GitlabInstance)
}

func NewGitlabInstance() (g *GitlabInstance) {
	return new(GitlabInstance)
}

func (g *GitlabInstance) AddRunner(newRunnerOptions *GitlabAddRunnerOptions) (createdRunner *GitlabRunner, err error) {
	if newRunnerOptions == nil {
		return nil, TracedError("newRunnerOptions is nil")
	}

	gitlabRunners, err := g.GetGitlabRunners()
	if err != nil {
		return nil, err
	}

	createdRunner, err = gitlabRunners.AddRunner(newRunnerOptions)
	if err != nil {
		return nil, err
	}

	return createdRunner, nil
}

func (g *GitlabInstance) Authenticate(authOptions *GitlabAuthenticationOptions) (err error) {
	if authOptions == nil {
		return TracedError("authOptions is nil")
	}

	fqdn, err := g.GetFqdn()
	if err != nil {
		return err
	}

	if authOptions.Verbose {
		LogInfof("Authenticate against gitlab '%s' started.", fqdn)
	}

	g.nativeClient = nil

	if authOptions.IsAccessTokenSet() {
		accessToken, err := authOptions.GetAccessToken()
		if err != nil {
			return err
		}

		apiV4Url, err := g.GetApiV4Url()
		if err != nil {
			return err
		}

		nativeClient, err := gitlab.NewClient(
			accessToken,
			gitlab.WithBaseURL(apiV4Url),
		)
		if err != nil {
			return TracedError(err.Error())
		}

		g.nativeClient = nativeClient
		g.currentlyUsedAccessToken = &accessToken
	}

	for _, gopassPath := range authOptions.AccessTokensFromGopass {
		credentialExists, err := Gopass().CredentialExists(gopassPath)
		if err != nil {
			return err
		}

		if !credentialExists {
			if authOptions.Verbose {
				LogInfof(
					"Gopass credential '%s' does not exist and can therefore not be used to authenticate against gitlab.",
					gopassPath,
				)
			}
			continue
		}

		getSecretOptions := NewGopassSecretOptions()
		getSecretOptions.SetGopassPath(gopassPath)
		accessToken, err := Gopass().GetCredentialValueAsString(getSecretOptions)
		if err != nil {
			return err
		}

		apiV4Url, err := g.GetApiV4Url()
		if err != nil {
			return err
		}

		nativeClient, err := gitlab.NewClient(
			accessToken,
			gitlab.WithBaseURL(apiV4Url),
		)
		if err != nil {
			return TracedError(err.Error())
		}

		g.nativeClient = nativeClient
		g.currentlyUsedAccessToken = &accessToken
	}

	if g.nativeClient == nil {
		return TracedErrorf("No authentication method for gitlab '%s' worked.", fqdn)
	}

	if authOptions.Verbose {
		LogInfof("Authenticate against gitlab '%s' finished.", fqdn)
	}

	return nil
}

func (g *GitlabInstance) CheckProjectByPathExists(projectPath string, verbose bool) (projectExists bool, err error) {
	if projectPath == "" {
		return false, TracedError("projectPath is empty string")
	}

	projectExists, err = g.ProjectByProjectPathExists(projectPath, verbose)
	if err != nil {
		return false, err
	}

	if !projectExists {
		errorMessage := fmt.Sprintf("Gitlab project '%s' does not exist.", projectPath)

		if verbose {
			LogError(errorMessage)
		}

		return false, TracedError(errorMessage)
	}

	return projectExists, nil
}

func (g *GitlabInstance) CheckRunnerStatusOk(runnerName string, verbose bool) (isStatusOk bool, err error) {
	if len(runnerName) <= 0 {
		return false, TracedError("runnerName is empty string")
	}

	gitlabRunners, err := g.GetGitlabRunners()
	if err != nil {
		return false, err
	}

	isStatusOk, err = gitlabRunners.CheckRunnerStatusOk(runnerName, verbose)
	if err != nil {
		return false, err
	}

	return isStatusOk, nil
}

func (g *GitlabInstance) CreateAccessToken(options *GitlabCreateAccessTokenOptions) (newToken string, err error) {
	if options == nil {
		return "", TracedError("options is nil")
	}

	users, err := g.GetGitlabUsers()
	if err != nil {
		return "", err
	}

	newToken, err = users.CreateAccessToken(options)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (g *GitlabInstance) CreateGroup(createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup, err error) {
	if createOptions == nil {
		return nil, TracedError("createOptions is nil")
	}

	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	createdGroup, err = gitlabGroups.CreateGroup(createOptions)
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}

func (g *GitlabInstance) CreateProject(createOptions *GitlabCreateProjectOptions) (gitlabProject *GitlabProject, err error) {
	if createOptions == nil {
		return nil, TracedError("createOptions is nil")
	}

	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	gitlabProject, err = gitlabProjects.CreateProject(createOptions)
	if err != nil {
		return nil, err
	}

	return gitlabProject, err
}

func (g *GitlabInstance) GetApiV4Url() (v4ApiUrl string, err error) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		return "", err
	}

	v4ApiUrl = fmt.Sprintf("https://%s/api/v4", fqdn)

	return v4ApiUrl, nil
}

func (g *GitlabInstance) GetCurrentlyUsedAccessToken() (gitlabAccessToken string, err error) {
	if g.currentlyUsedAccessToken == nil {
		return "", TracedError("currentlyUsedAccessToken not set")
	}

	return *g.currentlyUsedAccessToken, nil
}

func (g *GitlabInstance) GetDockerContainerOnGitlabHost(containerName string, sshUserName string) (dockerContainer *DockerContainer, err error) {
	if len(containerName) <= 0 {
		return nil, TracedError("containerName is empty string")
	}

	gitlabHost, err := g.GetHost()
	if err != nil {
		return nil, err
	}

	if len(sshUserName) > 0 {
		err = gitlabHost.SetSshUserName(sshUserName)
		if err != nil {
			return nil, err
		}
	}

	dockerContainer, err = gitlabHost.GetDockerContainerByName(containerName)
	if err != nil {
		return nil, err
	}

	return dockerContainer, nil
}

func (g *GitlabInstance) GetFqdn() (fqdn string, err error) {
	if len(g.fqdn) <= 0 {
		return "", TracedError("fqdn not set")
	}

	return g.fqdn, nil
}

func (g *GitlabInstance) GetGitlabGroups() (gitlabGroups *GitlabGroups, err error) {
	gitlabGroups = NewGitlabGroups()

	err = gitlabGroups.SetGitlab(g)
	if err != nil {
		return nil, err
	}

	return gitlabGroups, nil
}

func (g *GitlabInstance) GetGitlabProjectById(projectId int, verbose bool) (gitlabProject *GitlabProject, err error) {
	if projectId <= 0 {
		return nil, TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
	}

	gitlabProject = NewGitlabProject()
	err = gitlabProject.SetGitlab(g)
	if err != nil {
		return nil, err
	}

	err = gitlabProject.SetId(projectId)
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabInstance) GetGitlabProjectByPath(projectPath string, verbose bool) (gitlabProject *GitlabProject, err error) {
	if len(projectPath) <= 0 {
		return nil, TracedError("projectPath is empty string")
	}

	projectId, err := g.GetProjectIdByPath(projectPath, verbose)
	if err != nil {
		return nil, err
	}

	gitlabProject, err = g.GetGitlabProjectById(projectId, verbose)
	if err != nil {
		return nil, err
	}

	return gitlabProject, nil
}

func (g *GitlabInstance) GetGitlabProjects() (gitlabProjects *GitlabProjects, err error) {
	gitlabProjects = NewGitlabProjects()

	err = gitlabProjects.SetGitlab(g)
	if err != nil {
		return nil, err
	}

	return gitlabProjects, nil
}

func (g *GitlabInstance) GetGitlabRunners() (gitlabRunners *GitlabRunnersService, err error) {
	gitlabRunners = NewGitlabRunners()
	err = gitlabRunners.SetGitlab(g)
	if err != nil {
		return nil, err
	}

	return gitlabRunners, nil
}

func (g *GitlabInstance) GetGitlabSettings() (gitlabSettings *GitlabSettings, err error) {
	gitlabSettings = NewGitlabSettings()
	err = gitlabSettings.SetGitlab(g)
	if err != nil {
		return nil, err
	}

	return gitlabSettings, nil
}

func (g *GitlabInstance) GetGitlabUsers() (gitlabUsers *GitlabUsers, err error) {
	gitlabUsers = NewGitlabUsers()
	err = gitlabUsers.SetGitlab(g)
	if err != nil {
		return nil, err
	}

	return gitlabUsers, nil
}

func (g *GitlabInstance) GetHost() (gitlabHost *Host, err error) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		return nil, err
	}

	gitlabHost, err = GetHostByHostname(fqdn)
	if err != nil {
		return
	}

	return gitlabHost, nil
}

func (g *GitlabInstance) GetNativeClient() (nativeClient *gitlab.Client, err error) {
	if g.nativeClient == nil {
		return nil, TracedError("nativeClient not set")
	}

	return g.nativeClient, nil
}

func (g *GitlabInstance) GetNativeTagsService() (nativeTagsService *gitlab.TagsService, err error) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeTagsService = nativeClient.Tags
	if nativeTagsService == nil {
		return nil, TracedError("nativeTagsService is nil after evaluation")
	}

	return nativeTagsService, nil
}

func (g *GitlabInstance) GetPersonalAccessTokenList(verbose bool) (personalAccessTokens []*GitlabPersonalAccessToken, err error) {
	personalTokens, err := g.GetPersonalAccessTokens()
	if err != nil {
		return nil, err
	}

	personalAccessTokens, err = personalTokens.GetPersonalAccessTokenList(verbose)
	if err != nil {
		return nil, err
	}

	return personalAccessTokens, nil
}

func (g *GitlabInstance) GetPersonalAccessTokens() (tokens *GitlabPersonalAccessTokenService, err error) {
	tokens = NewGitlabPersonalAccessTokenService()

	err = tokens.SetGitlab(g)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (g *GitlabInstance) GetProjectIdByPath(projectPath string, verbose bool) (projectId int, err error) {
	if len(projectPath) <= 0 {
		return -1, TracedError("projectPath is empty string")
	}

	nativeClient, err := g.GetNativeClient()
	if err != nil {
		return -1, err
	}

	nativeProject, _, err := nativeClient.Projects.GetProject(projectPath, &gitlab.GetProjectOptions{})
	if err != nil {
		return -1, TracedError(err.Error())
	}

	projectId = nativeProject.ID
	if projectId <= 0 {
		return -1, TracedErrorf("Invalid project id returned by nativeProject: '%d'", projectId)
	}

	if verbose {
		LogInfof("Gitlab project '%s' has id '%d'", projectPath, projectId)
	}

	return projectId, nil
}

func (g *GitlabInstance) GetProjectPathList(verbose bool) (projectPaths []string, err error) {
	project, err := g.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	projectPaths, err = project.GetProjectPathList(verbose)
	if err != nil {
		return nil, err
	}

	return projectPaths, nil
}

func (g *GitlabInstance) GetRunnerByName(name string) (runner *GitlabRunner, err error) {
	if len(name) <= 0 {
		return nil, TracedError("name is empty string")
	}

	runners, err := g.GetGitlabRunners()
	if err != nil {
		return nil, err
	}

	runner, err = runners.GetRunnerByName(name)
	if err != nil {
		return nil, err
	}

	return runner, nil
}

func (g *GitlabInstance) GetUserByUsername(username string) (gitlabUser *GitlabUser, err error) {
	if len(username) <= 0 {
		return nil, TracedError("username is empty string")
	}

	gitlabUsers, err := g.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	gitlabUser, err = gitlabUsers.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return gitlabUser, nil
}

func (g *GitlabInstance) GetUserId() (userId int, err error) {
	users, err := g.GetGitlabUsers()
	if err != nil {
		return -1, err
	}

	userId, err = users.GetUserId()
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (g *GitlabInstance) GetUserNameList(verbose bool) (userNames []string, err error) {
	gitlabUsers, err := g.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	userNames, err = gitlabUsers.GetUserNames()
	if err != nil {
		return nil, err
	}

	return userNames, nil
}

func (g *GitlabInstance) GroupByGroupPathExists(groupPath string) (groupExists bool, err error) {
	if len(groupPath) <= 0 {
		return false, TracedError("groupPath is empty string")
	}

	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		return false, err
	}

	groupExists, err = gitlabGroups.GroupByGroupPathExists(groupPath)
	if err != nil {
		return false, err
	}

	return groupExists, nil
}

func (g *GitlabInstance) MustAddRunner(newRunnerOptions *GitlabAddRunnerOptions) (createdRunner *GitlabRunner) {
	createdRunner, err := g.AddRunner(newRunnerOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdRunner
}

func (g *GitlabInstance) MustAuthenticate(authOptions *GitlabAuthenticationOptions) {
	err := g.Authenticate(authOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustCheckProjectByPathExists(projectPath string, verbose bool) (projectExists bool) {
	projectExists, err := g.CheckProjectByPathExists(projectPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabInstance) MustCheckRunnerStatusOk(runnerName string, verbose bool) (isStatusOk bool) {
	isStatusOk, err := g.CheckRunnerStatusOk(runnerName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isStatusOk
}

func (g *GitlabInstance) MustCreateAccessToken(options *GitlabCreateAccessTokenOptions) (newToken string) {
	newToken, err := g.CreateAccessToken(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newToken
}

func (g *GitlabInstance) MustCreateGroup(createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup) {
	createdGroup, err := g.CreateGroup(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdGroup
}

func (g *GitlabInstance) MustCreateProject(createOptions *GitlabCreateProjectOptions) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.CreateProject(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabInstance) MustGetApiV4Url() (v4ApiUrl string) {
	v4ApiUrl, err := g.GetApiV4Url()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return v4ApiUrl
}

func (g *GitlabInstance) MustGetCurrentlyUsedAccessToken() (gitlabAccessToken string) {
	gitlabAccessToken, err := g.GetCurrentlyUsedAccessToken()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabAccessToken
}

func (g *GitlabInstance) MustGetDockerContainerOnGitlabHost(containerName string, sshUserName string) (dockerContainer *DockerContainer) {
	dockerContainer, err := g.GetDockerContainerOnGitlabHost(containerName, sshUserName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dockerContainer
}

func (g *GitlabInstance) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabInstance) MustGetGitlabGroups() (gitlabGroups *GitlabGroups) {
	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabGroups
}

func (g *GitlabInstance) MustGetGitlabProjectById(projectId int, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProjectById(projectId, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabInstance) MustGetGitlabProjectByPath(projectPath string, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProjectByPath(projectPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabInstance) MustGetGitlabProjects() (gitlabProjects *GitlabProjects) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProjects
}

func (g *GitlabInstance) MustGetGitlabRunners() (gitlabRunners *GitlabRunnersService) {
	gitlabRunners, err := g.GetGitlabRunners()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabRunners
}

func (g *GitlabInstance) MustGetGitlabSettings() (gitlabSettings *GitlabSettings) {
	gitlabSettings, err := g.GetGitlabSettings()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabSettings
}

func (g *GitlabInstance) MustGetGitlabUsers() (gitlabUsers *GitlabUsers) {
	gitlabUsers, err := g.GetGitlabUsers()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabUsers
}

func (g *GitlabInstance) MustGetHost() (gitlabHost *Host) {
	gitlabHost, err := g.GetHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabHost
}

func (g *GitlabInstance) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabInstance) MustGetNativeTagsService() (nativeTagsService *gitlab.TagsService) {
	nativeTagsService, err := g.GetNativeTagsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeTagsService
}

func (g *GitlabInstance) MustGetPersonalAccessTokenList(verbose bool) (personalAccessTokens []*GitlabPersonalAccessToken) {
	personalAccessTokens, err := g.GetPersonalAccessTokenList(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return personalAccessTokens
}

func (g *GitlabInstance) MustGetPersonalAccessTokens() (tokens *GitlabPersonalAccessTokenService) {
	tokens, err := g.GetPersonalAccessTokens()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tokens
}

func (g *GitlabInstance) MustGetProjectIdByPath(projectPath string, verbose bool) (projectId int) {
	projectId, err := g.GetProjectIdByPath(projectPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabInstance) MustGetProjectPathList(verbose bool) (projectPaths []string) {
	projectPaths, err := g.GetProjectPathList(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectPaths
}

func (g *GitlabInstance) MustGetRunnerByName(name string) (runner *GitlabRunner) {
	runner, err := g.GetRunnerByName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return runner
}

func (g *GitlabInstance) MustGetUserByUsername(username string) (gitlabUser *GitlabUser) {
	gitlabUser, err := g.GetUserByUsername(username)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabUser
}

func (g *GitlabInstance) MustGetUserId() (userId int) {
	userId, err := g.GetUserId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userId
}

func (g *GitlabInstance) MustGetUserNameList(verbose bool) (userNames []string) {
	userNames, err := g.GetUserNameList(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userNames
}

func (g *GitlabInstance) MustGroupByGroupPathExists(groupPath string) (groupExists bool) {
	groupExists, err := g.GroupByGroupPathExists(groupPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return groupExists
}

func (g *GitlabInstance) MustProjectByProjectIdExists(projectId int, verbose bool) (projectExists bool) {
	projectExists, err := g.ProjectByProjectIdExists(projectId, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabInstance) MustProjectByProjectPathExists(projectPath string, verbose bool) (projectExists bool) {
	projectExists, err := g.ProjectByProjectPathExists(projectPath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabInstance) MustRecreatePersonalAccessToken(createOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string) {
	newToken, err := g.RecreatePersonalAccessToken(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newToken
}

func (g *GitlabInstance) MustRemoveAllRunners(verbose bool) {
	err := g.RemoveAllRunners(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustResetAccessToken(options *GitlabResetAccessTokenOptions) {
	err := g.ResetAccessToken(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustResetUserPassword(resetOptions *GitlabResetPasswordOptions) {
	err := g.ResetUserPassword(resetOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustSetCurrentlyUsedAccessToken(currentlyUsedAccessToken *string) {
	err := g.SetCurrentlyUsedAccessToken(currentlyUsedAccessToken)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustSetFqdn(fqdn string) {
	err := g.SetFqdn(fqdn)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustSetNativeClient(nativeClient *gitlab.Client) {
	err := g.SetNativeClient(nativeClient)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustUseUnauthenticatedClient(verbose bool) {
	err := g.UseUnauthenticatedClient(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) ProjectByProjectIdExists(projectId int, verbose bool) (projectExists bool, err error) {
	if projectId <= 0 {
		return false, TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
	}

	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return false, err
	}

	projectExists, err = gitlabProjects.ProjectByProjectIdExists(projectId, verbose)
	if err != nil {
		return false, err
	}

	return projectExists, nil
}

func (g *GitlabInstance) ProjectByProjectPathExists(projectPath string, verbose bool) (projectExists bool, err error) {
	if len(projectPath) <= 0 {
		return false, TracedError("projectPath is empty string")
	}

	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return false, err
	}

	projectExists, err = gitlabProjects.ProjectByProjectPathExists(projectPath, verbose)
	if err != nil {
		return false, err
	}

	return projectExists, nil
}

func (g *GitlabInstance) RecreatePersonalAccessToken(createOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string, err error) {
	if createOptions == nil {
		return "", TracedError("createOptions is nil")
	}

	personalAccessTokens, err := g.GetPersonalAccessTokens()
	if err != nil {
		return "", err
	}

	newToken, err = personalAccessTokens.RecreateToken(createOptions)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (g *GitlabInstance) RemoveAllRunners(verbose bool) (err error) {
	runners, err := g.GetGitlabRunners()
	if err != nil {
		return err
	}

	err = runners.RemoveAllRunners(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabInstance) ResetAccessToken(options *GitlabResetAccessTokenOptions) (err error) {
	if options == nil {
		return TracedError("options is nil")
	}

	fqdn, err := g.GetFqdn()
	if err != nil {
		return err
	}

	username, err := options.GetUsername()
	if err != nil {
		return err
	}

	accessTokenName := "PERSONAL_ACCESS_TOKEN"

	if options.Verbose {
		LogInfof("Reset access token '%s' for user '%s' on gitlab '%s' started.", accessTokenName, username, fqdn)
	}

	gitlabContainer, err := g.GetDockerContainerOnGitlabHost(
		options.GitlabContainerNameOnGitlabHost,
		options.SshUserNameForGitlabHost,
	)
	if err != nil {
		return err
	}

	newToken, err := RandomGenerator().GetRandomString(30)
	if err != nil {
		return err
	}

	if options.Verbose {
		LogInfo("Going to create new token using gitlab-rails. Can take up to 30 seconds.")
	}

	expirationDate := time.Now().Add(time.Hour * 24 * 30)
	expirationDateString := expirationDate.Format("2006-01-02")
	_, err = gitlabContainer.RunCommand(
		&RunCommandOptions{
			Command: []string{
				"gitlab-rails",
				"runner",
				"token = User.find_by_username('" + username + "').personal_access_tokens.create(scopes: [:api], name: '" + accessTokenName + "', expires_at: '" + expirationDateString + "'); token.set_token('" + newToken + "'); token.save!",
			},
		},
	)
	if err != nil {
		return err
	}

	gopassOptions := NewGopassSecretOptions()
	gopassOptions.Verbose = options.Verbose
	gopassOptions.Overwrite = true
	gopassOptions.SetGopassPath(options.GopassPathToStoreNewToken)
	err = Gopass().InsertSecret(newToken, gopassOptions)
	if err != nil {
		return err
	}

	if options.Verbose {
		LogInfof("Reset access token '%s' for user '%s' on gitlab '%s' finished.", accessTokenName, username, fqdn)
	}

	return nil
}

func (g *GitlabInstance) ResetUserPassword(resetOptions *GitlabResetPasswordOptions) (err error) {

	if resetOptions == nil {
		return TracedError("resetOptions is nil")
	}

	fqdn, err := g.GetFqdn()
	if err != nil {
		return err
	}

	username, err := resetOptions.GetUsername()
	if err != nil {
		return err
	}

	if resetOptions.Verbose {
		LogInfof("Reset password for user '%s' on  gitlab '%s' started.", username, fqdn)
	}

	gitlabContainer, err := g.GetDockerContainerOnGitlabHost(
		resetOptions.GitlabContainerNameOnGitlabHost,
		resetOptions.SshUserNameForGitlabHost,
	)
	if err != nil {
		return err
	}

	newRootPassword, err := RandomGenerator().GetRandomString(12)
	if err != nil {
		return err
	}

	if resetOptions.Verbose {
		LogInfo("Going to reset password with gitlab-rake which usually takes several seconds.")
	}
	_, err = gitlabContainer.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"echo -ne '%s\n%s\n%s\n' | gitlab-rake \"gitlab:password:reset\"",
					username,
					newRootPassword,
					newRootPassword,
				),
			},
			Verbose: true,
		},
	)
	if err != nil {
		return err
	}

	gopassOptions := NewGopassSecretOptions()
	gopassOptions.Verbose = resetOptions.Verbose
	gopassOptions.Overwrite = true
	gopassOptions.SetGopassPath(resetOptions.GopassPathToStoreNewPassword)
	err = Gopass().InsertSecret(newRootPassword, gopassOptions)
	if err != nil {
		return err
	}
	gopassOptions.SecretBasename = "username"
	err = Gopass().InsertSecret(username, gopassOptions)
	if err != nil {
		return err
	}

	if resetOptions.Verbose {
		LogInfof("Reset password for user '%s' on  gitlab '%s' finished.", username, fqdn)
	}

	return nil

}

func (g *GitlabInstance) SetCurrentlyUsedAccessToken(currentlyUsedAccessToken *string) (err error) {
	if currentlyUsedAccessToken == nil {
		return TracedErrorf("currentlyUsedAccessToken is nil")
	}

	g.currentlyUsedAccessToken = currentlyUsedAccessToken

	return nil
}

func (g *GitlabInstance) SetFqdn(fqdn string) (err error) {
	if len(fqdn) <= 0 {
		return TracedError("fqdn is empty string")
	}

	fqdnUrl, err := GetUrlFromString(fqdn)
	if err != nil {
		return err
	}

	fqdnSane, err := fqdnUrl.GetFqdnAsString()
	if err != nil {
		return err
	}

	g.fqdn = fqdnSane

	return nil
}

func (g *GitlabInstance) SetNativeClient(nativeClient *gitlab.Client) (err error) {
	if nativeClient == nil {
		return TracedErrorf("nativeClient is nil")
	}

	g.nativeClient = nativeClient

	return nil
}

func (g *GitlabInstance) UseUnauthenticatedClient(verbose bool) (err error) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		return err
	}

	apiV4Url, err := g.GetApiV4Url()
	if err != nil {
		return err
	}

	nativeClient, err := gitlab.NewClient(
		"",
		gitlab.WithBaseURL(apiV4Url),
	)
	if err != nil {
		return TracedError(err.Error())
	}

	g.nativeClient = nativeClient

	if verbose {
		LogInfof("Unauthenticated gitlab client for '%s' is used.", fqdn)
	}

	return nil
}
