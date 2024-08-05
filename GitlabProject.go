package asciichgolangpublic

import (
	"github.com/xanzy/go-gitlab"
)

type GitlabProject struct {
	gitlab     *GitlabInstance
	id         int
	cachedPath string
}

func GetGitlabProjectByUrl(url *URL, authOptions []AuthenticationOption, verbose bool) (gitlabProject *GitlabProject, err error) {
	if url == nil {
		return nil, TracedErrorNil("url")
	}

	if authOptions == nil {
		return nil, TracedErrorNil("authOptions")
	}

	fqdnWithSheme, path, err := url.GetFqdnWitShemeAndPathAsString()
	if err != nil {
		return nil, err
	}

	gitlab, err := GetGitlabByFQDN(fqdnWithSheme)
	if err != nil {
		return nil, err
	}

	authOption, err := AuthenticationOptionsHandler().GetAuthenticationoptionsForServiceByUrl(authOptions, url)
	if err != nil {
		return nil, err
	}

	gitlabAuthenticationOption, ok := authOption.(*GitlabAuthenticationOptions)
	if !ok {
		return nil, TracedErrorf("Unable to get %v as GitlabAuthenticationOptions", authOption)
	}

	if authOptions != nil {
		err = gitlab.Authenticate(gitlabAuthenticationOption)
		if err != nil {
			return nil, err
		}
	}

	gitlabProject, err = gitlab.GetGitlabProjectByPath(path, verbose)
	if err != nil {
		return nil, err
	}

	return gitlabProject, err
}

func GetGitlabProjectByUrlFromString(urlString string, authOptions []AuthenticationOption, verbose bool) (gitlabProject *GitlabProject, err error) {
	if urlString == "" {
		return nil, TracedErrorEmptyString("urlString")
	}

	url, err := GetUrlFromString(urlString)
	if err != nil {
		return nil, err
	}

	return GetGitlabProjectByUrl(url, authOptions, verbose)
}

func MustGetGitlabProjectByUrl(url *URL, authOptions []AuthenticationOption, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := GetGitlabProjectByUrl(url, authOptions, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func MustGetGitlabProjectByUrlFromString(urlString string, authOptions []AuthenticationOption, verbose bool) (gitlabProject *GitlabProject) {
	gitlabProject, err := GetGitlabProjectByUrlFromString(urlString, authOptions, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func NewGitlabProject() (gitlabProject *GitlabProject) {
	return new(GitlabProject)
}

func (g *GitlabProject) Exists(verbose bool) (projectExists bool, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return false, err
	}

	projectId, err := g.GetId()
	if err != nil {
		return false, err
	}

	projectExists, err = gitlab.ProjectByProjectIdExists(projectId, verbose)
	if err != nil {
		return false, err
	}

	return projectExists, nil
}

func (g *GitlabProject) GetNewestVersion(verbose bool) (newestVersion Version, err error) {
	availableVersions, err := g.GetVersions(verbose)
	if err != nil {
		return nil, err
	}

	if len(availableVersions) <= 0 {
		return nil, TracedError("No versionTags returned")
	}

	newestVersion, err = Versions().GetLatestVersionFromSlice(availableVersions)
	if err != nil {
		return nil, err
	}

	return newestVersion, nil
}

func (g *GitlabProject) GetNewestVersionAsString(verbose bool) (newestVersionString string, err error) {
	newestVersion, err := g.GetNewestVersion(verbose)
	if err != nil {
		return "", err
	}

	newestVersionString, err = newestVersion.GetAsString()
	if err != nil {
		return "", err
	}

	return newestVersionString, nil
}

func (g *GitlabProject) GetPath() (projectPath string, err error) {
	gitlabProjects, err := g.GetGitlabProjects()
	if err != nil {
		return "", err
	}

	projectsService, err := gitlabProjects.GetNativeProjectsService()
	if err != nil {
		return "", err
	}

	projectId, err := g.GetId()
	if err != nil {
		return "", err
	}

	nativeProject, _, err := projectsService.GetProject(projectId, nil, nil)
	if err != nil {
		return "", TracedErrorf("Unable to get native project: '%w'", err)
	}

	projectPath = nativeProject.PathWithNamespace
	if projectPath == "" {
		return "", TracedError("projectPath is empty string after evaluation")
	}

	err = g.SetCachedPath(projectPath)
	if err != nil {
		return "", err
	}

	return projectPath, nil
}

func (g *GitlabProject) GetTags() (gitlabTags *GitlabTags, err error) {
	gitlabTags = NewGitlabTags()

	err = gitlabTags.SetGitlabProject(g)
	if err != nil {
		return nil, err
	}

	return gitlabTags, err
}

func (g *GitlabProject) GetVersionTagNames(verbose bool) (versionTagNames []string, err error) {
	tags, err := g.GetTags()
	if err != nil {
		return nil, err
	}

	versionTagNames, err = tags.GetVersionTagNames(verbose)
	if err != nil {
		return nil, err
	}

	return versionTagNames, nil
}

func (g *GitlabProject) GetVersionTags(verbose bool) (versionTags []*GitlabTag, err error) {
	tags, err := g.GetTags()
	if err != nil {
		return nil, err
	}

	versionTags, err = tags.GetVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	return versionTags, nil
}

func (g *GitlabProject) GetVersions(verbose bool) (versions []Version, err error) {
	versionTags, err := g.GetVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	versions = []Version{}

	for _, tag := range versionTags {
		versionName, err := tag.GetName()
		if err != nil {
			return nil, err
		}

		toAdd, err := Versions().GetNewVersionByString(versionName)
		if err != nil {
			return nil, err
		}

		versions = append(versions, toAdd)
	}

	return versions, nil
}

func (g *GitlabProject) IsCachedPathSet() (isSet bool) {
	return g.cachedPath != ""
}

func (g *GitlabProject) MustDeployKeyByNameExists(keyName string) (exists bool) {
	exists, err := g.DeployKeyByNameExists(keyName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabProject) MustExists(verbose bool) (projectExists bool) {
	projectExists, err := g.Exists(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectExists
}

func (g *GitlabProject) MustGetCachedPath() (path string) {
	path, err := g.GetCachedPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return path
}

func (g *GitlabProject) MustGetDeployKeyByName(keyName string) (projectDeployKey *GitlabProjectDeployKey) {
	projectDeployKey, err := g.GetDeployKeyByName(keyName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectDeployKey
}

func (g *GitlabProject) MustGetDeployKeys() (deployKeys *GitlabProjectDeployKeys) {
	deployKeys, err := g.GetDeployKeys()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return deployKeys
}

func (g *GitlabProject) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabProject) MustGetGitlabProjectDeployKeys() (projectDeployKeys *GitlabProjectDeployKeys) {
	projectDeployKeys, err := g.GetGitlabProjectDeployKeys()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectDeployKeys
}

func (g *GitlabProject) MustGetGitlabProjects() (projects *GitlabProjects) {
	projects, err := g.GetGitlabProjects()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projects
}

func (g *GitlabProject) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabProject) MustGetNativeProjectsService() (nativeGitlabProject *gitlab.ProjectsService) {
	nativeGitlabProject, err := g.GetNativeProjectsService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeGitlabProject
}

func (g *GitlabProject) MustGetNewestVersion(verbose bool) (newestVersion Version) {
	newestVersion, err := g.GetNewestVersion(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newestVersion
}

func (g *GitlabProject) MustGetNewestVersionAsString(verbose bool) (newestVersionString string) {
	newestVersionString, err := g.GetNewestVersionAsString(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newestVersionString
}

func (g *GitlabProject) MustGetPath() (projectPath string) {
	projectPath, err := g.GetPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectPath
}

func (g *GitlabProject) MustGetTags() (gitlabTags *GitlabTags) {
	gitlabTags, err := g.GetTags()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabTags
}

func (g *GitlabProject) MustGetVersionTagNames(verbose bool) (versionTagNames []string) {
	versionTagNames, err := g.GetVersionTagNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionTagNames
}

func (g *GitlabProject) MustGetVersionTags(verbose bool) (versionTags []*GitlabTag) {
	versionTags, err := g.GetVersionTags(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionTags
}

func (g *GitlabProject) MustGetVersions(verbose bool) (versions []Version) {
	versions, err := g.GetVersions(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versions
}

func (g *GitlabProject) MustMakePrivate(verbose bool) {
	err := g.MakePrivate(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustMakePublic(verbose bool) {
	err := g.MakePublic(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustRecreateDeployKey(keyOptions *GitlabCreateDeployKeyOptions) {
	err := g.RecreateDeployKey(keyOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustSetCachedPath(pathToCache string) {
	err := g.SetCachedPath(pathToCache)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabProject) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *GitlabProject) DeployKeyByNameExists(keyName string) (exists bool, err error) {
	if len(keyName) <= 0 {
		return false, TracedError("keyName is empty string")
	}

	deployKeys, err := p.GetDeployKeys()
	if err != nil {
		return false, err
	}

	exists, err = deployKeys.DeployKeyByNameExists(keyName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (p *GitlabProject) GetCachedPath() (path string, err error) {
	if !p.IsCachedPathSet() {
		_, err := p.GetPath()
		if err != nil {
			return "", err
		}
	}

	if len(p.cachedPath) <= 0 {
		return "", TracedError("cachedPath not set")
	}

	return p.cachedPath, nil
}

func (p *GitlabProject) GetDeployKeyByName(keyName string) (projectDeployKey *GitlabProjectDeployKey, err error) {
	if len(keyName) <= 0 {
		return nil, TracedError("keyName is nil")
	}

	deployKeys, err := p.GetGitlabProjectDeployKeys()
	if err != nil {
		return nil, err
	}

	projectDeployKey, err = deployKeys.GetGitlabProjectDeployKeyByName(keyName)
	if err != nil {
		return nil, err
	}

	return projectDeployKey, nil
}

func (p *GitlabProject) GetDeployKeys() (deployKeys *GitlabProjectDeployKeys, err error) {
	deployKeys = NewGitlabProjectDeployKeys()
	err = deployKeys.SetGitlabProject(p)
	if err != nil {
		return nil, err
	}

	return deployKeys, nil
}

func (p *GitlabProject) GetGitlab() (gitlab *GitlabInstance, err error) {
	if p.gitlab == nil {
		return nil, TracedError("gitlab is not set")
	}

	return p.gitlab, nil
}

func (p *GitlabProject) GetGitlabProjectDeployKeys() (projectDeployKeys *GitlabProjectDeployKeys, err error) {
	projectDeployKeys = NewGitlabProjectDeployKeys()

	err = projectDeployKeys.SetGitlabProject(p)
	if err != nil {
		return nil, err
	}

	return projectDeployKeys, nil
}

func (p *GitlabProject) GetGitlabProjects() (projects *GitlabProjects, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	projects, err = gitlab.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *GitlabProject) GetId() (id int, err error) {
	if p.id <= 0 {
		return -1, TracedError("id not set")
	}

	return p.id, nil
}

func (p *GitlabProject) GetNativeProjectsService() (nativeGitlabProject *gitlab.ProjectsService, err error) {
	projects, err := p.GetGitlabProjects()
	if err != nil {
		return nil, err
	}

	nativeGitlabProject, err = projects.GetNativeProjectsService()
	if err != nil {
		return nil, err
	}

	return nativeGitlabProject, nil
}

func (p *GitlabProject) MakePrivate(verbose bool) (err error) {
	nativeProjectsService, err := p.GetNativeProjectsService()
	if err != nil {
		return err
	}

	projectId, err := p.GetId()
	if err != nil {
		return err
	}

	var visibility = gitlab.PrivateVisibility

	_, _, err = nativeProjectsService.EditProject(
		projectId,
		&gitlab.EditProjectOptions{
			Visibility: &visibility,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof("Gitlab project '%v' made private.", projectId)
	}

	return nil
}

func (p *GitlabProject) MakePublic(verbose bool) (err error) {
	nativeProjectsService, err := p.GetNativeProjectsService()
	if err != nil {
		return err
	}

	projectId, err := p.GetId()
	if err != nil {
		return err
	}

	var visibility = gitlab.PublicVisibility

	_, _, err = nativeProjectsService.EditProject(
		projectId,
		&gitlab.EditProjectOptions{
			Visibility: &visibility,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof("Gitlab project '%v' made public.", projectId)
	}

	return nil
}

func (p *GitlabProject) RecreateDeployKey(keyOptions *GitlabCreateDeployKeyOptions) (err error) {
	if keyOptions == nil {
		return TracedError("keyOptions is nil")
	}

	deployKey, err := p.GetDeployKeyByName(keyOptions.Name)
	if err != nil {
		return err
	}

	err = deployKey.RecreateDeployKey(keyOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *GitlabProject) SetCachedPath(pathToCache string) (err error) {
	if len(pathToCache) <= 0 {
		return TracedError("pathToCache is empty string")
	}

	p.cachedPath = pathToCache

	return nil
}

func (p *GitlabProject) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}

func (p *GitlabProject) SetId(id int) (err error) {
	if id <= 0 {
		return TracedErrorf("invalid id = '%d'", id)
	}

	p.id = id

	return nil
}
