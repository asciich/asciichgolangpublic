package asciichgolangpublic

import (
	"slices"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabReleaseLinks struct {
	gitlabRelease *GitlabRelease
}

func NewGitlabReleaseLinks() (g *GitlabReleaseLinks) {
	return new(GitlabReleaseLinks)
}

func (g *GitlabReleaseLinks) CreateReleaseLink(createOptions *GitlabCreateReleaseLinkOptions) (createdReleaseLink *GitlabReleaseLink, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedErrorNil("createOptions")
	}

	nativeClient, err := g.GetNativeReleaseLinksClient()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl()
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

	exists, err := g.ReleaseLinkByNameExists(linkName, createOptions.Verbose)
	if err != nil {
		return nil, err
	}

	if exists {
		if createOptions.Verbose {
			logging.LogInfof(
				"Release link '%s' for release '%s' in gitlab project %s already exists. Skip creation.",
				linkName,
				tagName,
				projectUrl,
			)
		}
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

		if createOptions.Verbose {
			logging.LogChangedf(
				"Created release link '%s' in gitlab project %s .",
				linkUrl,
				projectUrl,
			)
		}
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

func (g *GitlabReleaseLinks) GetProjectId() (projectId int, err error) {
	release, err := g.GetGitlabRelease()
	if err != nil {
		return -1, err
	}

	projectId, err = release.GetProjectId()
	if err != nil {
		return -1, err
	}

	return projectId, nil
}

func (g *GitlabReleaseLinks) GetProjectIdAndUrl() (projectId int, projectUrl string, err error) {
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

func (g *GitlabReleaseLinks) GetProjectUrl() (projectUrl string, err error) {
	release, err := g.GetGitlabRelease()
	if err != nil {
		return "", err
	}

	projectUrl, err = release.GetProjectUrl()
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

func (g *GitlabReleaseLinks) HasReleaseLinks(verbose bool) (hasReleaseLinks bool, err error) {
	releaseLinks, err := g.ListReleaseLinks(verbose)
	if err != nil {
		return false, err
	}

	hasReleaseLinks = len(releaseLinks) > 0

	return hasReleaseLinks, nil
}

func (g *GitlabReleaseLinks) ListReleaseLinkNames(verbose bool) (releaseLinkNames []string, err error) {
	releaseLinks, err := g.ListReleaseLinks(verbose)
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

func (g *GitlabReleaseLinks) ListReleaseLinkUrls(verbose bool) (releaseLinkUrls []string, err error) {
	releaseLinks, err := g.ListReleaseLinks(verbose)
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

func (g *GitlabReleaseLinks) ListReleaseLinks(verbose bool) (releaseLinks []*GitlabReleaseLink, err error) {
	nativeClient, err := g.GetNativeReleaseLinksClient()
	if err != nil {
		return nil, err
	}

	projectId, projectUrl, err := g.GetProjectIdAndUrl()
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

func (g *GitlabReleaseLinks) MustCreateReleaseLink(createOptions *GitlabCreateReleaseLinkOptions) (createdReleaseLink *GitlabReleaseLink) {
	createdReleaseLink, err := g.CreateReleaseLink(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdReleaseLink
}

func (g *GitlabReleaseLinks) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabReleaseLinks) MustGetGitlabRelease() (gitlabRelease *GitlabRelease) {
	gitlabRelease, err := g.GetGitlabRelease()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabRelease
}

func (g *GitlabReleaseLinks) MustGetNativeReleaseLinksClient() (nativeClient *gitlab.ReleaseLinksService) {
	nativeClient, err := g.GetNativeReleaseLinksClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabReleaseLinks) MustGetProjectId() (projectId int) {
	projectId, err := g.GetProjectId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId
}

func (g *GitlabReleaseLinks) MustGetProjectIdAndUrl() (projectId int, projectUrl string) {
	projectId, projectUrl, err := g.GetProjectIdAndUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectId, projectUrl
}

func (g *GitlabReleaseLinks) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabReleaseLinks) MustGetReleaseLinkByName(linkName string) (releaseLink *GitlabReleaseLink) {
	releaseLink, err := g.GetReleaseLinkByName(linkName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return releaseLink
}

func (g *GitlabReleaseLinks) MustGetReleaseName() (releaseName string) {
	releaseName, err := g.GetReleaseName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return releaseName
}

func (g *GitlabReleaseLinks) MustHasReleaseLinks(verbose bool) (hasReleaseLinks bool) {
	hasReleaseLinks, err := g.HasReleaseLinks(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hasReleaseLinks
}

func (g *GitlabReleaseLinks) MustListReleaseLinkNames(verbose bool) (releaseLinkNames []string) {
	releaseLinkNames, err := g.ListReleaseLinkNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return releaseLinkNames
}

func (g *GitlabReleaseLinks) MustListReleaseLinkUrls(verbose bool) (releaseLinkUrls []string) {
	releaseLinkUrls, err := g.ListReleaseLinkUrls(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return releaseLinkUrls
}

func (g *GitlabReleaseLinks) MustListReleaseLinks(verbose bool) (releaseLinks []*GitlabReleaseLink) {
	releaseLinks, err := g.ListReleaseLinks(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return releaseLinks
}

func (g *GitlabReleaseLinks) MustReleaseLinkByNameExists(linkName string, verbose bool) (exists bool) {
	exists, err := g.ReleaseLinkByNameExists(linkName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabReleaseLinks) MustSetGitlabRelease(gitlabRelease *GitlabRelease) {
	err := g.SetGitlabRelease(gitlabRelease)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabReleaseLinks) ReleaseLinkByNameExists(linkName string, verbose bool) (exists bool, err error) {
	if linkName == "" {
		return false, tracederrors.TracedErrorEmptyString("linkName")
	}

	releaseNames, err := g.ListReleaseLinkNames(verbose)
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
