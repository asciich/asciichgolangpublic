package asciichgolangpublic

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabInstance struct {
	fqdn                     string
	nativeClient             *gitlab.Client
	currentlyUsedAccessToken *string
}

func GetGitlabByFQDN(fqdn string) (gitlab *GitlabInstance, err error) {
	if len(fqdn) <= 0 {
		return nil, tracederrors.TracedError("fqdn is empty string")
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
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func MustGetGitlabByFqdn(fqdn string) (gitlab *GitlabInstance) {
	gitlab, err := GetGitlabByFQDN(fqdn)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func NewGitlab() (gitlab *GitlabInstance) {
	return new(GitlabInstance)
}

func NewGitlabInstance() (g *GitlabInstance) {
	return new(GitlabInstance)
}

// Get the path to the personal projects which is "users/USERNAME/projects".
func (g *GitlabInstance) GetPersonalProjectsPath(verbose bool) (personalProjetsPath string, err error) {
	userName, err := g.GetCurrentUsersName(verbose)
	if err != nil {
		return "", err
	}

	personalProjetsPath = fmt.Sprintf("users/%s/projects", userName)

	return personalProjetsPath, nil
}

// Return the gitlab user name.
// This is the technical user name used by Gitlab.
//
// To get the human readable user name use `GetCurrentUsersName`.
func (g *GitlabInstance) GetCurrentUsersUsername(verbose bool) (currentUserName string, err error) {
	user, err := g.GetCurrentUser(verbose)
	if err != nil {
		return "", err
	}

	currentUserName, err = user.GetCachedUsername()
	if err != nil {
		return "", err
	}

	return currentUserName, nil
}

// Returns the `userId` of the currently logged in user.
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

// Returns the human readable gitlab user name also known as display name.
//
// For the technical user name use `GetCurrentUsersUsername`.
func (g *GitlabInstance) GetCurrentUsersName(verbose bool) (currentUserName string, err error) {
	user, err := g.GetCurrentUser(verbose)
	if err != nil {
		return "", err
	}

	currentUserName, err = user.GetCachedName()
	if err != nil {
		return "", err
	}

	return currentUserName, nil
}

func (g *GitlabInstance) AddRunner(newRunnerOptions *GitlabAddRunnerOptions) (createdRunner *GitlabRunner, err error) {
	if newRunnerOptions == nil {
		return nil, tracederrors.TracedError("newRunnerOptions is nil")
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
		return tracederrors.TracedError("authOptions is nil")
	}

	fqdn, err := g.GetFqdn()
	if err != nil {
		return err
	}

	if authOptions.Verbose {
		logging.LogInfof("Authenticate against gitlab '%s' started.", fqdn)
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
			return tracederrors.TracedError(err.Error())
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
				logging.LogInfof(
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
			return tracederrors.TracedError(err.Error())
		}

		g.nativeClient = nativeClient
		g.currentlyUsedAccessToken = &accessToken
	}

	if g.nativeClient == nil {
		return tracederrors.TracedErrorf("No authentication method for gitlab '%s' worked.", fqdn)
	}

	if authOptions.Verbose {
		logging.LogInfof("Authenticate against gitlab '%s' finished.", fqdn)
	}

	return nil
}

func (g *GitlabInstance) CheckProjectByPathExists(projectPath string, verbose bool) (projectExists bool, err error) {
	if projectPath == "" {
		return false, tracederrors.TracedError("projectPath is empty string")
	}

	projectExists, err = g.ProjectByProjectPathExists(projectPath, verbose)
	if err != nil {
		return false, err
	}

	if !projectExists {
		errorMessage := fmt.Sprintf("Gitlab project '%s' does not exist.", projectPath)

		if verbose {
			logging.LogError(errorMessage)
		}

		return false, tracederrors.TracedError(errorMessage)
	}

	return projectExists, nil
}

func (g *GitlabInstance) CheckRunnerStatusOk(runnerName string, verbose bool) (isStatusOk bool, err error) {
	if len(runnerName) <= 0 {
		return false, tracederrors.TracedError("runnerName is empty string")
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
		return "", tracederrors.TracedError("options is nil")
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

func (g *GitlabInstance) CreateGroupByPath(groupPath string, createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup, err error) {
	if groupPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("groupPath")
	}

	if createOptions == nil {
		return nil, tracederrors.TracedErrorNil("createOptions")
	}

	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	createdGroup, err = gitlabGroups.CreateGroup(groupPath, createOptions)
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}

func (g *GitlabInstance) CreatePersonalProject(projectName string, verbose bool) (personalProject *GitlabProject, err error) {
	if projectName == "" {
		return nil, tracederrors.TracedErrorEmptyString("projectName")
	}

	personalProject, err = g.GetPersonalProjectByName(projectName, verbose)
	if err != nil {
		return nil, err
	}

	err = personalProject.Create(verbose)
	if err != nil {
		return nil, err
	}

	return personalProject, nil
}

func (g *GitlabInstance) CreateProject(createOptions *GitlabCreateProjectOptions) (gitlabProject *GitlabProject, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
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

func (g *GitlabInstance) DeleteGroupByPath(groupPath string, verbose bool) (err error) {
	if groupPath == "" {
		return tracederrors.TracedErrorEmptyString("groupPath")
	}

	group, err := g.GetGroupByPath(groupPath, verbose)
	if err != nil {
		return err
	}

	err = group.Delete(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitlabInstance) GetApiV4Url() (v4ApiUrl string, err error) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		return "", err
	}

	v4ApiUrl = fmt.Sprintf("https://%s/api/v4", fqdn)

	return v4ApiUrl, nil
}

func (g *GitlabInstance) GetCurrentUser(verbose bool) (currentUser *GitlabUser, err error) {
	users, err := g.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	currentUser, err = users.GetUser()
	if err != nil {
		return nil, err
	}

	return currentUser, nil
}

func (g *GitlabInstance) GetCurrentlyUsedAccessToken() (gitlabAccessToken string, err error) {
	if g.currentlyUsedAccessToken == nil {
		return "", tracederrors.TracedError("currentlyUsedAccessToken not set")
	}

	return *g.currentlyUsedAccessToken, nil
}

func (g *GitlabInstance) GetDeepCopy() (copy *GitlabInstance) {
	copy = NewGitlab()

	*copy = *g

	if g.currentlyUsedAccessToken != nil {
		copy.currentlyUsedAccessToken = g.currentlyUsedAccessToken
	}

	return copy
}

func (g *GitlabInstance) GetFqdn() (fqdn string, err error) {
	if len(g.fqdn) <= 0 {
		return "", tracederrors.TracedError("fqdn not set")
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
		return nil, tracederrors.TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
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
		return nil, tracederrors.TracedError("projectPath is empty string")
	}

	exists, err := g.ProjectByProjectPathExists(projectPath, verbose)
	if err != nil {
		return nil, err
	}

	if exists {
		projectId, err := g.GetProjectIdByPath(projectPath, verbose)
		if err != nil {
			return nil, err
		}

		gitlabProject, err = g.GetGitlabProjectById(projectId, verbose)
		if err != nil {
			return nil, err
		}

		err = gitlabProject.SetCachedPath(projectPath)
		if err != nil {
			return nil, err
		}
	} else {
		gitlabProject = NewGitlabProject()
		err = gitlabProject.SetCachedPath(projectPath)
		if err != nil {
			return nil, err
		}

		err = gitlabProject.SetGitlab(g)
		if err != nil {
			return nil, err
		}
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

func (g *GitlabInstance) GetGroupById(id int, verbose bool) (gitlabGroup *GitlabGroup, err error) {
	groups, err := g.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	gitlabGroup, err = groups.GetGroupById(id, verbose)
	if err != nil {
		return nil, err
	}

	return gitlabGroup, nil
}

func (g *GitlabInstance) GetGroupByPath(groupPath string, verbose bool) (gitlabGroup *GitlabGroup, err error) {
	groups, err := g.GetGitlabGroups()
	if err != nil {
		return nil, err
	}

	gitlabGroup, err = groups.GetGroupByPath(groupPath, verbose)
	if err != nil {
		return nil, err
	}

	return gitlabGroup, nil
}

func (g *GitlabInstance) GetNativeBranchesClient() (nativeClient *gitlab.BranchesService, err error) {
	client, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeClient = client.Branches
	if nativeClient == nil {
		return nil, tracederrors.TracedError("nativeClient is nil after evaluation")
	}

	return nativeClient, nil
}

func (g *GitlabInstance) GetNativeClient() (nativeClient *gitlab.Client, err error) {
	if g.nativeClient == nil {
		return nil, tracederrors.TracedError("nativeClient not set")
	}

	return g.nativeClient, nil
}

func (g *GitlabInstance) GetNativeMergeRequestsService() (nativeClient *gitlab.MergeRequestsService, err error) {
	client, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeClient = client.MergeRequests
	if nativeClient == nil {
		return nil, tracederrors.TracedError("nativeClient is nil after evaluation")
	}

	return nativeClient, nil
}

func (g *GitlabInstance) GetNativeReleaseLinksClient() (nativeClient *gitlab.ReleaseLinksService, err error) {
	client, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeClient = client.ReleaseLinks

	if nativeClient == nil {
		return nil, tracederrors.TracedError("native client is nil after evaluation.")
	}

	return nativeClient, nil
}

func (g *GitlabInstance) GetNativeReleasesClient() (nativeReleasesClient *gitlab.ReleasesService, err error) {
	gitlabClient, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeReleasesClient = gitlabClient.Releases

	if nativeReleasesClient == nil {
		return nil, tracederrors.TracedError(
			"Native releases client is empty string after evaluation.",
		)
	}

	return nativeReleasesClient, nil
}

func (g *GitlabInstance) GetNativeRepositoriesClient() (nativeRepositoriesClient *gitlab.RepositoriesService, err error) {
	client, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeRepositoriesClient = client.Repositories
	if nativeRepositoriesClient == nil {
		return nil, tracederrors.TracedError("Repositories is nil after evaluation")
	}

	return nativeRepositoriesClient, nil
}

func (g *GitlabInstance) GetNativeRepositoryFilesClient() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService, err error) {
	client, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeRepositoryFilesClient = client.RepositoryFiles
	if nativeRepositoryFilesClient == nil {
		return nil, tracederrors.TracedError("nativeRepositoryFilesClient is nil after evaluation")
	}

	return nativeRepositoryFilesClient, nil
}

func (g *GitlabInstance) GetNativeTagsService() (nativeTagsService *gitlab.TagsService, err error) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeTagsService = nativeClient.Tags
	if nativeTagsService == nil {
		return nil, tracederrors.TracedError("nativeTagsService is nil after evaluation")
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

func (g *GitlabInstance) GetPersonalProjectByName(projectName string, verbose bool) (project *GitlabProject, err error) {
	if projectName == "" {
		return nil, tracederrors.TracedErrorEmptyString("projectName")
	}

	personalProjectsPath, err := g.GetPersonalProjectsPath(verbose)
	if err != nil {
		return nil, err
	}

	personalProjectsPath = astrings.EnsureSuffix(personalProjectsPath, "/")

	projectPath := astrings.EnsurePrefix(projectName, personalProjectsPath)

	project, err = g.GetGitlabProjectByPath(projectPath, verbose)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (g *GitlabInstance) GetProjectIdByPath(projectPath string, verbose bool) (projectId int, err error) {
	if len(projectPath) <= 0 {
		return -1, tracederrors.TracedError("projectPath is empty string")
	}

	projects, err := g.GetGitlabProjects()
	if err != nil {
		return -1, err
	}

	projectId, err = projects.GetProjectIdByProjectPath(projectPath, verbose)
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabInstance) GetProjectPathList(options *GitlabgetProjectListOptions) (projectPaths []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	project, err := g.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	projectPaths, err = project.GetProjectPathList(options)
	if err != nil {
		return nil, err
	}

	return projectPaths, nil
}

func (g *GitlabInstance) GetRunnerByName(name string) (runner *GitlabRunner, err error) {
	if len(name) <= 0 {
		return nil, tracederrors.TracedError("name is empty string")
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

func (g *GitlabInstance) GetUserById(id int) (gitlabUser *GitlabUser, err error) {
	if id <= 0 {
		return nil, tracederrors.TracedErrorf("id '%d' is invalid", id)
	}

	gitlabUsers, err := g.GetGitlabUsers()
	if err != nil {
		return nil, err
	}

	gitlabUser, err = gitlabUsers.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return gitlabUser, nil
}

func (g *GitlabInstance) GetUserByUsername(username string) (gitlabUser *GitlabUser, err error) {
	if len(username) <= 0 {
		return nil, tracederrors.TracedError("username is empty string")
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

func (g *GitlabInstance) GroupByGroupPathExists(groupPath string, verbose bool) (groupExists bool, err error) {
	if len(groupPath) <= 0 {
		return false, tracederrors.TracedError("groupPath is empty string")
	}

	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		return false, err
	}

	groupExists, err = gitlabGroups.GroupByGroupPathExists(groupPath, verbose)
	if err != nil {
		return false, err
	}

	return groupExists, nil
}

func (g *GitlabInstance) MustAddRunner(newRunnerOptions *GitlabAddRunnerOptions) (createdRunner *GitlabRunner) {
	createdRunner, err := g.AddRunner(newRunnerOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRunner
}

func (g *GitlabInstance) MustAuthenticate(authOptions *GitlabAuthenticationOptions) {
	err := g.Authenticate(authOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustCheckProjectByPathExists(projectPath string, verbose bool) (projectExists bool) {
	projectExists, err := g.CheckProjectByPathExists(projectPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabInstance) MustCheckRunnerStatusOk(runnerName string, verbose bool) (isStatusOk bool) {
	isStatusOk, err := g.CheckRunnerStatusOk(runnerName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isStatusOk
}

func (g *GitlabInstance) MustCreateGroupByPath(groupPath string, createOptions *GitlabCreateGroupOptions) (createdGroup *GitlabGroup) {
	createdGroup, err := g.CreateGroupByPath(groupPath, createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdGroup
}

func (g *GitlabInstance) MustCreatePersonalProject(projectName string, verbose bool) (personalProject *GitlabProject) {
	personalProject, err := g.CreatePersonalProject(projectName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return personalProject
}

func (g *GitlabInstance) MustCreateProject(createOptions *GitlabCreateProjectOptions) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.CreateProject(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabInstance) MustDeleteGroupByPath(groupPath string, verbose bool) {
	err := g.DeleteGroupByPath(groupPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustGetApiV4Url() (v4ApiUrl string) {
	v4ApiUrl, err := g.GetApiV4Url()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return v4ApiUrl
}

func (g *GitlabInstance) MustGetCurrentUser(verbose bool) (currentUser *GitlabUser) {
	currentUser, err := g.GetCurrentUser(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return currentUser
}

func (g *GitlabInstance) MustGetCurrentUserName(verbose bool) (currentUserName string) {
	currentUserName, err := g.GetCurrentUsersName(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return currentUserName
}

func (g *GitlabInstance) MustGetCurrentUsersName(verbose bool) (currentUserName string) {
	currentUserName, err := g.GetCurrentUsersName(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return currentUserName
}

func (g *GitlabInstance) MustGetCurrentUsersUsername(verbose bool) (currentUserName string) {
	currentUserName, err := g.GetCurrentUsersUsername(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return currentUserName
}

func (g *GitlabInstance) MustGetCurrentlyUsedAccessToken() (gitlabAccessToken string) {
	gitlabAccessToken, err := g.GetCurrentlyUsedAccessToken()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabAccessToken
}

func (g *GitlabInstance) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabInstance) MustGetGitlabGroups() (gitlabGroups *GitlabGroups) {
	gitlabGroups, err := g.GetGitlabGroups()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabGroups
}

func (g *GitlabInstance) MustGetGitlabProjectById(projectId int, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProjectById(projectId, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabInstance) MustGetGitlabProjectByPath(projectPath string, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProjectByPath(projectPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabInstance) MustGetGitlabProjects() (gitlabProjects *GitlabProjects) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabProjects
}

func (g *GitlabInstance) MustGetGitlabRunners() (gitlabRunners *GitlabRunnersService) {
	gitlabRunners, err := g.GetGitlabRunners()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabRunners
}

func (g *GitlabInstance) MustGetGitlabSettings() (gitlabSettings *GitlabSettings) {
	gitlabSettings, err := g.GetGitlabSettings()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabSettings
}

func (g *GitlabInstance) MustGetGitlabUsers() (gitlabUsers *GitlabUsers) {
	gitlabUsers, err := g.GetGitlabUsers()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabUsers
}

func (g *GitlabInstance) MustGetGroupById(id int, verbose bool) (gitlabGroup *GitlabGroup) {
	gitlabGroup, err := g.GetGroupById(id, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabGroup
}

func (g *GitlabInstance) MustGetGroupByPath(groupPath string, verbose bool) (gitlabGroup *GitlabGroup) {
	gitlabGroup, err := g.GetGroupByPath(groupPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabGroup
}

func (g *GitlabInstance) MustGetNativeBranchesClient() (nativeClient *gitlab.BranchesService) {
	nativeClient, err := g.GetNativeBranchesClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabInstance) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabInstance) MustGetNativeMergeRequestsService() (nativeClient *gitlab.MergeRequestsService) {
	nativeClient, err := g.GetNativeMergeRequestsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabInstance) MustGetNativeReleaseLinksClient() (nativeClient *gitlab.ReleaseLinksService) {
	nativeClient, err := g.GetNativeReleaseLinksClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabInstance) MustGetNativeReleasesClient() (nativeReleasesClient *gitlab.ReleasesService) {
	nativeReleasesClient, err := g.GetNativeReleasesClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeReleasesClient
}

func (g *GitlabInstance) MustGetNativeRepositoriesClient() (nativeRepositoriesClient *gitlab.RepositoriesService) {
	nativeRepositoriesClient, err := g.GetNativeRepositoriesClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeRepositoriesClient
}

func (g *GitlabInstance) MustGetNativeRepositoryFilesClient() (nativeRepositoryFilesClient *gitlab.RepositoryFilesService) {
	nativeRepositoryFilesClient, err := g.GetNativeRepositoryFilesClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeRepositoryFilesClient
}

func (g *GitlabInstance) MustGetNativeTagsService() (nativeTagsService *gitlab.TagsService) {
	nativeTagsService, err := g.GetNativeTagsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeTagsService
}

func (g *GitlabInstance) MustGetPersonalAccessTokenList(verbose bool) (personalAccessTokens []*GitlabPersonalAccessToken) {
	personalAccessTokens, err := g.GetPersonalAccessTokenList(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return personalAccessTokens
}

func (g *GitlabInstance) MustGetPersonalAccessTokens() (tokens *GitlabPersonalAccessTokenService) {
	tokens, err := g.GetPersonalAccessTokens()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tokens
}

func (g *GitlabInstance) MustGetPersonalProjectByName(projectName string, verbose bool) (project *GitlabProject) {
	project, err := g.GetPersonalProjectByName(projectName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return project
}

func (g *GitlabInstance) MustGetPersonalProjectsPath(verbose bool) (personalProjetsPath string) {
	personalProjetsPath, err := g.GetPersonalProjectsPath(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return personalProjetsPath
}

func (g *GitlabInstance) MustGetProjectIdByPath(projectPath string, verbose bool) (projectId int) {
	projectId, err := g.GetProjectIdByPath(projectPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabInstance) MustGetProjectPathList(options *GitlabgetProjectListOptions) (projectPaths []string) {
	projectPaths, err := g.GetProjectPathList(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectPaths
}

func (g *GitlabInstance) MustGetRunnerByName(name string) (runner *GitlabRunner) {
	runner, err := g.GetRunnerByName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return runner
}

func (g *GitlabInstance) MustGetUserById(id int) (gitlabUser *GitlabUser) {
	gitlabUser, err := g.GetUserById(id)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabUser
}

func (g *GitlabInstance) MustGetUserByUsername(username string) (gitlabUser *GitlabUser) {
	gitlabUser, err := g.GetUserByUsername(username)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabUser
}

func (g *GitlabInstance) MustGetUserId() (userId int) {
	userId, err := g.GetUserId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userId
}

func (g *GitlabInstance) MustGetUserNameList(verbose bool) (userNames []string) {
	userNames, err := g.GetUserNameList(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userNames
}

func (g *GitlabInstance) MustGroupByGroupPathExists(groupPath string, verbose bool) (groupExists bool) {
	groupExists, err := g.GroupByGroupPathExists(groupPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return groupExists
}

func (g *GitlabInstance) MustProjectByProjectIdExists(projectId int, verbose bool) (projectExists bool) {
	projectExists, err := g.ProjectByProjectIdExists(projectId, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabInstance) MustProjectByProjectPathExists(projectPath string, verbose bool) (projectExists bool) {
	projectExists, err := g.ProjectByProjectPathExists(projectPath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabInstance) MustRecreatePersonalAccessToken(createOptions *GitlabCreatePersonalAccessTokenOptions) (newToken string) {
	newToken, err := g.RecreatePersonalAccessToken(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newToken
}

func (g *GitlabInstance) MustRemoveAllRunners(verbose bool) {
	err := g.RemoveAllRunners(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustResetAccessToken(options *GitlabResetAccessTokenOptions) {
	err := g.ResetAccessToken(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustResetUserPassword(resetOptions *GitlabResetPasswordOptions) {
	err := g.ResetUserPassword(resetOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustSetCurrentlyUsedAccessToken(currentlyUsedAccessToken *string) {
	err := g.SetCurrentlyUsedAccessToken(currentlyUsedAccessToken)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustSetFqdn(fqdn string) {
	err := g.SetFqdn(fqdn)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustSetNativeClient(nativeClient *gitlab.Client) {
	err := g.SetNativeClient(nativeClient)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) MustUseUnauthenticatedClient(verbose bool) {
	err := g.UseUnauthenticatedClient(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabInstance) ProjectByProjectIdExists(projectId int, verbose bool) (projectExists bool, err error) {
	if projectId <= 0 {
		return false, tracederrors.TracedErrorf("projectId '%d' <= 0 is invalid", projectId)
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
		return false, tracederrors.TracedError("projectPath is empty string")
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
		return "", tracederrors.TracedError("createOptions is nil")
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
		return tracederrors.TracedError("options is nil")
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
		logging.LogInfof("Reset access token '%s' for user '%s' on gitlab '%s' started.", accessTokenName, username, fqdn)
	}

	return tracederrors.TracedErrorNotImplemented()
	/*

		gitlabContainer, err := docker.GetDockerContainerOnHost(
			g,
			options.GitlabContainerNameOnGitlabHost,
		)
		if err != nil {
			return err
		}
		newToken, err := RandomGenerator().GetRandomString(30)
		if err != nil {
			return err
		}

		if options.Verbose {
			logging.LogInfo("Going to create new token using gitlab-rails. Can take up to 30 seconds.")
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
			logging.LogInfof("Reset access token '%s' for user '%s' on gitlab '%s' finished.", accessTokenName, username, fqdn)
		}

		return nil
	*/
}

func (g *GitlabInstance) ResetUserPassword(resetOptions *GitlabResetPasswordOptions) (err error) {

	if resetOptions == nil {
		return tracederrors.TracedError("resetOptions is nil")
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
		logging.LogInfof("Reset password for user '%s' on  gitlab '%s' started.", username, fqdn)
	}

	return tracederrors.TracedErrorNotImplemented()
	/*

		gitlabContainer, err := docker.GetDockerContainerOnHost(
			g,
			resetOptions.GitlabContainerNameOnGitlabHost,
		)
		if err != nil {
			return err
		}

		newRootPassword, err := RandomGenerator().GetRandomString(12)
		if err != nil {
			return err
		}

		if resetOptions.Verbose {
			logging.LogInfo("Going to reset password with gitlab-rake which usually takes several seconds.")
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
			logging.LogInfof("Reset password for user '%s' on  gitlab '%s' finished.", username, fqdn)
		}

		return nil

	*/
}

func (g *GitlabInstance) SetCurrentlyUsedAccessToken(currentlyUsedAccessToken *string) (err error) {
	if currentlyUsedAccessToken == nil {
		return tracederrors.TracedErrorf("currentlyUsedAccessToken is nil")
	}

	g.currentlyUsedAccessToken = currentlyUsedAccessToken

	return nil
}

func (g *GitlabInstance) SetFqdn(fqdn string) (err error) {
	if len(fqdn) <= 0 {
		return tracederrors.TracedError("fqdn is empty string")
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
		return tracederrors.TracedErrorf("nativeClient is nil")
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
		return tracederrors.TracedError(err.Error())
	}

	g.nativeClient = nativeClient

	if verbose {
		logging.LogInfof("Unauthenticated gitlab client for '%s' is used.", fqdn)
	}

	return nil
}
