package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabReleases struct {
	gitlabProject *GitlabProject
}

func NewGitlabReleases() (g *GitlabReleases) {
	return new(GitlabReleases)
}

func (g *GitlabReleases) CreateRelease(ctx context.Context, createReleaseOptions *GitlabCreateReleaseOptions) (createdRelease *GitlabRelease, err error) {
	if createReleaseOptions == nil {
		return nil, tracederrors.TracedErrorNil("createReleaseOptions")
	}

	nativeClient, err := g.GetNativeReleasesClient()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl(ctx)
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

	exists, err := g.ReleaseByNameExists(ctx, releaseName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Release '%s' already exists in gitlab project %s . Skip creation.", releaseName, projectUrl)
	} else {
		_, _, err = nativeClient.CreateRelease(
			projectId,
			&gitlab.CreateReleaseOptions{
				Name:        &releaseName,
				TagName:     &releaseName,
				Description: &description,
			},
		)
		if err != nil {
			return nil, tracederrors.TracedErrorf(
				"Create release '%s' in gitlab project %s failed: %w",
				releaseName,
				projectUrl,
				err,
			)
		}

		logging.LogChangedByCtxf(ctx, "Created release '%s' in gitlab project %s", releaseName, projectUrl)
	}

	createdRelease, err = g.GetGitlabReleaseByName(releaseName)
	if err != nil {
		return nil, err
	}

	return createdRelease, nil
}

func (g *GitlabReleases) DeleteAllReleases(ctx context.Context, options *GitlabDeleteReleaseOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	releaseList, err := g.ListReleases(ctx)
	if err != nil {
		return err
	}

	for _, toDelete := range releaseList {
		err = toDelete.Delete(ctx, options)
		if err != nil {
			return err
		}
	}

	projectUrl, err := g.GetProjectUrl(ctx)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Deleted '%d' releases from gitlab project %s .", len(releaseList), projectUrl)

	return err
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
		return nil, tracederrors.TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabReleases) GetGitlabReleaseByName(releaseName string) (gitlabRelease *GitlabRelease, err error) {
	if releaseName == "" {
		return nil, tracederrors.TracedErrorEmptyString("releaseName")
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

func (g *GitlabReleases) GetProjectId(ctx context.Context) (projectId int, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return -1, err
	}

	projectId, err = project.GetId(ctx)
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabReleases) GetProjectIdAndUrl(ctx context.Context) (projectId int, projectUrl string, err error) {
	projectId, err = g.GetProjectId(ctx)
	if err != nil {
		return -1, "", err
	}

	projectUrl, err = g.GetProjectUrl(ctx)
	if err != nil {
		return -1, "", err
	}

	return projectId, projectUrl, err
}

func (g *GitlabReleases) GetProjectUrl(ctx context.Context) (projectUrl string, err error) {
	project, err := g.GetGitlabProject()
	if err != nil {
		return "", err
	}

	projectUrl, err = project.GetProjectUrl(ctx)
	if err != nil {
		return "", err
	}

	return projectUrl, nil
}

func (g *GitlabReleases) ListReleases(ctx context.Context) (releaseList []*GitlabRelease, err error) {
	projectId, projectUrl, err := g.GetProjectIdAndUrl(ctx)
	if err != nil {
		return nil, err
	}

	nativeClient, err := g.GetNativeReleasesClient()
	if err != nil {
		return nil, err
	}

	rawReleases, _, err := nativeClient.ListReleases(
		projectId,
		&gitlab.ListReleasesOptions{},
		nil,
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to list releases of gitlab project %s : %w",
			projectUrl,
			err,
		)
	}

	releaseList = []*GitlabRelease{}
	for _, raw := range rawReleases {
		toAdd, err := g.GetGitlabReleaseByName(raw.Name)
		if err != nil {
			return nil, err
		}

		releaseList = append(releaseList, toAdd)
	}

	return releaseList, nil
}

func (g *GitlabReleases) ReleaseByNameExists(ctx context.Context, releaseName string) (exists bool, err error) {
	if releaseName == "" {
		return false, tracederrors.TracedErrorEmptyString("releaseName")
	}

	release, err := g.GetGitlabReleaseByName(releaseName)
	if err != nil {
		return false, err
	}

	exists, err = release.Exists(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (g *GitlabReleases) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return tracederrors.TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}
