package asciichgolangpublic

import (
	"context"
	"slices"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabReleaseLinks struct {
	gitlabRelease *GitlabRelease
}

func NewGitlabReleaseLinks() (g *GitlabReleaseLinks) {
	return new(GitlabReleaseLinks)
}

func (g *GitlabReleaseLinks) CreateReleaseLink(ctx context.Context, createOptions *GitlabCreateReleaseLinkOptions) (createdReleaseLink *GitlabReleaseLink, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedErrorNil("createOptions")
	}

	nativeClient, err := g.GetNativeReleaseLinksClient()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl(ctx)
	if err != nil {
		return nil, err
	}

	tagName, err := g.GetReleaseName()
	if err != nil {
		return nil, err
	}

	linkName, linkUrl, err := createOptions.GetNameAndUrl()
	if err != nil {
		return nil, err
	}

	exists, err := g.ReleaseLinkByNameExists(ctx, linkName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Release link '%s' for release '%s' in gitlab project %s already exists. Skip creation.", linkName, tagName, projectUrl)
	} else {
		_, _, err = nativeClient.CreateReleaseLink(
			projectId,
			tagName,
			&gitlab.CreateReleaseLinkOptions{
				Name: &linkName,
				URL:  &linkUrl,
			},
		)
		if err != nil {
			return nil, tracederrors.TracedErrorf(
				"Create release link '%s' in gitlab project %s failed: %w",
				linkName,
				projectUrl,
				err,
			)
		}

		logging.LogChangedByCtxf(ctx, "Created release link '%s' in gitlab project %s .", linkUrl, projectUrl)
	}

	createdReleaseLink, err = g.GetReleaseLinkByName(linkName)
	if err != nil {
		return nil, err
	}

	return createdReleaseLink, nil
}

func (g *GitlabReleaseLinks) GetGitlab() (gitlab *GitlabInstance, err error) {
	release, err := g.GetGitlabRelease()
	if err != nil {
		return nil, err
	}

	gitlab, err = release.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, err
}

func (g *GitlabReleaseLinks) GetGitlabRelease() (gitlabRelease *GitlabRelease, err error) {
	if g.gitlabRelease == nil {
		return nil, tracederrors.TracedErrorf("gitlabRelease not set")
	}

	return g.gitlabRelease, nil
}

func (g *GitlabReleaseLinks) GetNativeReleaseLinksClient() (nativeClient *gitlab.ReleaseLinksService, err error) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeReleaseLinksClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (g *GitlabReleaseLinks) GetProjectId(ctx context.Context) (projectId int, err error) {
	release, err := g.GetGitlabRelease()
	if err != nil {
		return -1, err
	}

	projectId, err = release.GetProjectId(ctx)
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabReleaseLinks) GetProjectIdAndUrl(ctx context.Context) (projectId int, projectUrl string, err error) {
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

func (g *GitlabReleaseLinks) GetProjectUrl(ctx context.Context) (projectUrl string, err error) {
	release, err := g.GetGitlabRelease()
	if err != nil {
		return "", err
	}

	projectUrl, err = release.GetProjectUrl(ctx)
	if err != nil {
		return "", err
	}

	return projectUrl, nil
}

func (g *GitlabReleaseLinks) GetReleaseLinkByName(linkName string) (releaseLink *GitlabReleaseLink, err error) {
	if linkName == "" {
		return nil, tracederrors.TracedErrorEmptyString("linkName")
	}

	releaseLink = NewGitlabReleaseLink()

	err = releaseLink.SetGitlabReleaseLinks(g)
	if err != nil {
		return nil, err
	}

	err = releaseLink.SetName(linkName)
	if err != nil {
		return nil, err
	}

	return releaseLink, nil
}

func (g *GitlabReleaseLinks) GetReleaseName() (releaseName string, err error) {
	release, err := g.GetGitlabRelease()
	if err != nil {
		return "", err
	}

	releaseName, err = release.GetName()
	if err != nil {
		return "", err
	}

	return releaseName, nil
}

func (g *GitlabReleaseLinks) HasReleaseLinks(ctx context.Context) (hasReleaseLinks bool, err error) {
	releaseLinks, err := g.ListReleaseLinks(ctx)
	if err != nil {
		return false, err
	}

	hasReleaseLinks = len(releaseLinks) > 0

	return hasReleaseLinks, nil
}

func (g *GitlabReleaseLinks) ListReleaseLinkNames(ctx context.Context) (releaseLinkNames []string, err error) {
	releaseLinks, err := g.ListReleaseLinks(ctx)
	if err != nil {
		return nil, err
	}

	releaseLinkNames = []string{}
	for _, link := range releaseLinks {
		toAdd, err := link.GetName()
		if err != nil {
			return nil, err
		}

		releaseLinkNames = append(releaseLinkNames, toAdd)
	}

	return releaseLinkNames, nil
}

func (g *GitlabReleaseLinks) ListReleaseLinkUrls(ctx context.Context) (releaseLinkUrls []string, err error) {
	releaseLinks, err := g.ListReleaseLinks(ctx)
	if err != nil {
		return nil, err
	}

	releaseLinkUrls = []string{}
	for _, link := range releaseLinks {
		toAdd, err := link.GetCachedUrl()
		if err != nil {
			return nil, err
		}

		releaseLinkUrls = append(releaseLinkUrls, toAdd)
	}

	return releaseLinkUrls, nil
}

func (g *GitlabReleaseLinks) ListReleaseLinks(ctx context.Context) (releaseLinks []*GitlabReleaseLink, err error) {
	nativeClient, err := g.GetNativeReleaseLinksClient()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl(ctx)
	if err != nil {
		return nil, err
	}

	tagName, err := g.GetReleaseName()
	if err != nil {
		return nil, err
	}

	rawReleaseLinks, _, err := nativeClient.ListReleaseLinks(projectId, tagName, nil)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Failed to list release links for release '%s' in gitlab project %s : %w",
			tagName,
			projectUrl,
			err,
		)
	}

	releaseLinks = []*GitlabReleaseLink{}
	for _, raw := range rawReleaseLinks {
		toAdd, err := g.GetReleaseLinkByName(raw.Name)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetCachedUrl(raw.URL)
		if err != nil {
			return nil, err
		}

		releaseLinks = append(releaseLinks, toAdd)
	}

	return releaseLinks, nil
}

func (g *GitlabReleaseLinks) ReleaseLinkByNameExists(ctx context.Context, linkName string) (exists bool, err error) {
	if linkName == "" {
		return false, tracederrors.TracedErrorEmptyString("linkName")
	}

	releaseNames, err := g.ListReleaseLinkNames(ctx)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(releaseNames, linkName)

	return exists, nil
}

func (g *GitlabReleaseLinks) SetGitlabRelease(gitlabRelease *GitlabRelease) (err error) {
	if gitlabRelease == nil {
		return tracederrors.TracedErrorf("gitlabRelease is nil")
	}

	g.gitlabRelease = gitlabRelease

	return nil
}
