package asciichgolangpublic

import "github.com/xanzy/go-gitlab"

type GitlabReleases struct {
	gitlabProject *GitlabProject
}

func NewGitlabReleases() (g *GitlabReleases) {
	return new(GitlabReleases)
}

func (g *GitlabReleases) CreateRelease(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease, err error) {
	if createReleaseOptions == nil {
		return nil, TracedErrorNil("createReleaseOptions")
	}

	nativeClient, err := g.GetNativeReleasesClient()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		return nil, err
	}

	releaseName, err := createReleaseOptions.GetName()
	if err != nil {
		return nil, err
	}

	description, err := createReleaseOptions.GetDescription()
	if err != nil {
		return nil, err
	}

	_, _, err = nativeClient.CreateRelease(
		projectId,
		&gitlab.CreateReleaseOptions{
			Name:        &releaseName,
			TagName:     &releaseName,
			Description: &description,
		},
	)
	if err != nil {
		return nil, TracedErrorf(
			"Create release '%s' in gitlab project %s failed: %w",
			releaseName,
			projectUrl,
			err,
		)
	}

	createdRelease, err = g.GetGitlabReleaseByName(releaseName)
	if err != nil {
		return nil, err
	}

	if createReleaseOptions.Verbose {
		LogChangedf(
			"Created release '%s' in gitlab project %s",
			releaseName,
			projectUrl,
		)
	}

	return createdRelease, nil
}

func (g *GitlabReleases) GetGitlab() (gitlab *GitlabInstance, err error) {
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

func (g *GitlabReleases) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabReleases) GetGitlabReleaseByName(releaseName string) (gitlabRelease *GitlabRelease, err error) {
	if releaseName == "" {
		return nil, TracedErrorEmptyString("releaseName")
	}

	gitlabRelease = NewGitlabRelease()

	err = gitlabRelease.SetGitlabReleases(g)
	if err != nil {
		return nil, err
	}

	err = gitlabRelease.SetName(releaseName)
	if err != nil {
		return nil, err
	}

	return gitlabRelease, nil
}

func (g *GitlabReleases) GetNativeReleasesClient() (nativeReleasesClient *gitlab.ReleasesService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeReleasesClient, err = gitlab.GetNativeReleasesClient()
	if err != nil {
		return nil, err
	}

	return nativeReleasesClient, nil
}

func (g *GitlabReleases) GetProjectId() (projectId int, err error) {
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

func (g *GitlabReleases) GetProjectIdAndUrl() (projectId int, projectUrl string, err error) {
	projectId, err = g.GetProjectId()
	if err != nil {
		return -1, "", err
	}

	projectUrl, err = g.GetProjectUrl()
	if err != nil {
		return -1, "", err
	}

	return projectId, projectUrl, err
}

func (g *GitlabReleases) GetProjectUrl() (projectUrl string, err error) {
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

func (g *GitlabReleases) MustCreateRelease(createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease) {
	createdRelease, err := g.CreateRelease(createReleaseOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdRelease
}

func (g *GitlabReleases) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabReleases) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabReleases) MustGetGitlabReleaseByName(releaseName string) (gitlabRelease *GitlabRelease) {
	gitlabRelease, err := g.GetGitlabReleaseByName(releaseName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabRelease
}

func (g *GitlabReleases) MustGetNativeReleasesClient() (nativeReleasesClient *gitlab.ReleasesService) {
	nativeReleasesClient, err := g.GetNativeReleasesClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeReleasesClient
}

func (g *GitlabReleases) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabReleases) MustGetProjectIdAndUrl() (projectId int, projectUrl string) {
	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectId, projectUrl
}

func (g *GitlabReleases) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabReleases) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabReleases) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}
