package asciichgolangpublic

import (
	"errors"

	"github.com/xanzy/go-gitlab"
)

var ErrGitlabReleaseNotFound = errors.New("gitlab release not found")

type GitlabRelease struct {
	name           string
	gitlabReleases *GitlabReleases
}

func NewGitlabRelease() (g *GitlabRelease) {
	return new(GitlabRelease)
}

func (g *GitlabRelease) CreateReleaseLink(createOptions *GitlabCreateReleaseLinkOptions) (createdReleaseLink *GitlabReleaseLink, err error) {
	releaseLinks, err := g.GetGitlabReleaseLinks()
	if err != nil {
		return nil, err
	}

	createdReleaseLink, err = releaseLinks.CreateReleaseLink(createOptions)
	if err != nil {
		return nil, err
	}

	return createdReleaseLink, nil
}

func (g *GitlabRelease) Delete(deleteOptions *GitlabDeleteReleaseOptions) (err error) {
	if deleteOptions == nil {
		return TracedErrorNil("deleteOptions")
	}

	exists, err := g.Exists(deleteOptions.Verbose)
	if err != nil {
		return err
	}

	projectUrl, releaseName, err := g.GetProjectUrlAndReleaseName()
	if err != nil {
		return err
	}

	if exists {
		projectId, err := g.GetProjectId()
		if err != nil {
			return err
		}

		nativeClient, err := g.GetNativeReleasesClient()
		if err != nil {
			return err
		}

		_, _, err = nativeClient.DeleteRelease(
			projectId,
			releaseName,
			nil,
		)
		if err != nil {
			return err
		}

		LogChangedf(
			"Release '%s' on gitlab project '%s' deleted.",
			projectUrl,
			releaseName,
		)
	} else {
		LogInfof(
			"Release '%s' on gitlab project '%s' is already absent. Skip delete.",
			projectUrl,
			releaseName,
		)
	}

	deleteCorrespondingTag := deleteOptions.GetDeleteCorrespondingTag()

	if deleteCorrespondingTag {
		err = g.DeleteCorrespondingTag(
			deleteOptions.Verbose,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GitlabRelease) DeleteCorrespondingTag(verbose bool) (err error) {
	name, err := g.GetName()
	if err != nil {
		return err
	}

	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Delete tag corresonding to release '%s' in gitlab project %s started.",
			name,
			projectUrl,
		)
	}

	tag, err := g.GetTag()
	if err != nil {
		return err
	}

	err = tag.Delete(verbose)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Delete tag corresonding to release '%s' in gitlab project %s finished.",
			name,
			projectUrl,
		)
	}

	return nil
}

func (g *GitlabRelease) Exists(verbose bool) (exists bool, err error) {
	projectUrl, releaseName, err := g.GetProjectUrlAndReleaseName()
	if err != nil {
		return false, err
	}

	exists = true
	_, err = g.GetRawResponse()
	if err != nil {
		if errors.Is(err, ErrGitlabReleaseNotFound) {
			exists = false
		} else {
			return false, err
		}
	}

	if exists {
		LogInfof(
			"Gitlab Release '%s' exists in project %s .",
			releaseName,
			projectUrl,
		)
	} else {
		LogInfof(
			"Gitlab Release '%s' does not exist in project %s .",
			releaseName,
			projectUrl,
		)
	}

	return exists, nil
}

func (g *GitlabRelease) GetGitlab() (gitlab *GitlabInstance, err error) {
	releases, err := g.GetGitlabReleases()
	if err != nil {
		return nil, err
	}

	gitlab, err = releases.GetGitlab()
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}

func (g *GitlabRelease) GetGitlabProject() (project *GitlabProject, err error) {
	releases, err := g.GetGitlabReleases()
	if err != nil {
		return nil, err
	}

	project, err = releases.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (g *GitlabRelease) GetGitlabReleaseLinks() (gitlabReleaseLinks *GitlabReleaseLinks, err error) {
	gitlabReleaseLinks = NewGitlabReleaseLinks()

	err = gitlabReleaseLinks.SetGitlabRelease(g)
	if err != nil {
		return nil, err
	}
	return gitlabReleaseLinks, nil
}

func (g *GitlabRelease) GetGitlabReleases() (gitlabReleases *GitlabReleases, err error) {
	if g.gitlabReleases == nil {
		return nil, TracedErrorf("gitlabReleases not set")
	}

	return g.gitlabReleases, nil
}

func (g *GitlabRelease) GetName() (name string, err error) {
	if g.name == "" {
		return "", TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GitlabRelease) GetNativeReleasesClient() (nativeReleasesClient *gitlab.ReleasesService, err error) {
	releases, err := g.GetGitlabReleases()
	if err != nil {
		return nil, err
	}

	nativeReleasesClient, err = releases.GetNativeReleasesClient()
	if err != nil {
		return nil, err
	}

	return nativeReleasesClient, nil
}

func (g *GitlabRelease) GetProjectId() (pid int, err error) {
	releases, err := g.GetGitlabReleases()
	if err != nil {
		return -1, err
	}

	pid, err = releases.GetProjectId()
	if err != nil {
		return -1, err
	}

	return pid, nil
}

func (g *GitlabRelease) GetProjectUrl() (projectUrl string, err error) {
	releases, err := g.GetGitlabReleases()
	if err != nil {
		return "", err
	}

	projectUrl, err = releases.GetProjectUrl()
	if err != nil {
		return "", err
	}

	return projectUrl, nil
}

func (g *GitlabRelease) GetProjectUrlAndReleaseName() (projectUrl string, releaseName string, err error) {
	projectUrl, err = g.GetProjectUrl()
	if err != nil {
		return "", "", err
	}

	releaseName, err = g.GetName()
	if err != nil {
		return "", "", err
	}

	return projectUrl, releaseName, nil
}

func (g *GitlabRelease) GetRawResponse() (rawRelease *gitlab.Release, err error) {
	nativeClient, err := g.GetNativeReleasesClient()
	if err != nil {
		return nil, err
	}

	projectId, err := g.GetProjectId()
	if err != nil {
		return nil, err
	}

	name, err := g.GetName()
	if err != nil {
		return nil, err
	}

	projectUrl, releaseName, err := g.GetProjectUrlAndReleaseName()
	if err != nil {
		return nil, err
	}

	rawRelease, _, err = nativeClient.GetRelease(projectId, name)
	if err != nil {
		if err.Error() == "404 Not Found" {
			return nil, TracedErrorf(
				"%w, Project %s release '%s'",
				ErrGitlabReleaseNotFound,
				projectUrl,
				releaseName,
			)
		}

		return nil, TracedErrorf(
			"Failed to GetRawResponse for gitlab release '%s' for project %s : '%w'",
			releaseName,
			projectUrl,
			err,
		)
	}

	if rawRelease == nil {
		return nil, TracedError("rawRelease is empty string after evaluation")
	}

	return rawRelease, nil
}

func (g *GitlabRelease) GetTag() (tag *GitlabTag, err error) {
	name, err := g.GetName()
	if err != nil {
		return nil, err
	}

	project, err := g.GetGitlabProject()
	if err != nil {
		return nil, err
	}

	tag, err = project.GetTagByName(name)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (g *GitlabRelease) HasReleaseLinks(verbose bool) (hasReleaseLinks bool, err error) {
	releaseLinks, err := g.GetGitlabReleaseLinks()
	if err != nil {
		return false, err
	}

	hasReleaseLinks, err = releaseLinks.HasReleaseLinks(verbose)
	if err != nil {
		return false, err
	}

	return hasReleaseLinks, nil
}

func (g *GitlabRelease) ListReleaseLinkUrls(verbose bool) (releaseLinkUrls []string, err error) {
	releaseLinks, err := g.GetGitlabReleaseLinks()
	if err != nil {
		return nil, err
	}

	releaseLinkUrls, err = releaseLinks.ListReleaseLinkUrls(verbose)
	if err != nil {
		return nil, err
	}

	return releaseLinkUrls, nil
}

func (g *GitlabRelease) MustCreateReleaseLink(createOptions *GitlabCreateReleaseLinkOptions) (createdReleaseLink *GitlabReleaseLink) {
	createdReleaseLink, err := g.CreateReleaseLink(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdReleaseLink
}

func (g *GitlabRelease) MustDelete(deleteOptions *GitlabDeleteReleaseOptions) {
	err := g.Delete(deleteOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRelease) MustDeleteCorrespondingTag(verbose bool) {
	err := g.DeleteCorrespondingTag(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRelease) MustExists(verbose bool) (exists bool) {
	exists, err := g.Exists(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabRelease) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabRelease) MustGetGitlabProject() (project *GitlabProject) {
	project, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return project
}

func (g *GitlabRelease) MustGetGitlabReleaseLinks() (gitlabReleaseLinks *GitlabReleaseLinks) {
	gitlabReleaseLinks, err := g.GetGitlabReleaseLinks()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabReleaseLinks
}

func (g *GitlabRelease) MustGetGitlabReleases() (gitlabReleases *GitlabReleases) {
	gitlabReleases, err := g.GetGitlabReleases()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabReleases
}

func (g *GitlabRelease) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabRelease) MustGetNativeReleasesClient() (nativeReleasesClient *gitlab.ReleasesService) {
	nativeReleasesClient, err := g.GetNativeReleasesClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeReleasesClient
}

func (g *GitlabRelease) MustGetProjectId() (pid int) {
	pid, err := g.GetProjectId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return pid
}

func (g *GitlabRelease) MustGetProjectUrl() (projectUrl string) {
	projectUrl, err := g.GetProjectUrl()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectUrl
}

func (g *GitlabRelease) MustGetProjectUrlAndReleaseName() (projectUrl string, releaseName string) {
	projectUrl, releaseName, err := g.GetProjectUrlAndReleaseName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return projectUrl, releaseName
}

func (g *GitlabRelease) MustGetRawResponse() (rawRelease *gitlab.Release) {
	rawRelease, err := g.GetRawResponse()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rawRelease
}

func (g *GitlabRelease) MustGetTag() (tag *GitlabTag) {
	tag, err := g.GetTag()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tag
}

func (g *GitlabRelease) MustHasReleaseLinks(verbose bool) (hasReleaseLinks bool) {
	hasReleaseLinks, err := g.HasReleaseLinks(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasReleaseLinks
}

func (g *GitlabRelease) MustListReleaseLinkUrls(verbose bool) (releaseLinkUrls []string) {
	releaseLinkUrls, err := g.ListReleaseLinkUrls(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return releaseLinkUrls
}

func (g *GitlabRelease) MustSetGitlabReleases(gitlabReleases *GitlabReleases) {
	err := g.SetGitlabReleases(gitlabReleases)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRelease) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRelease) SetGitlabReleases(gitlabReleases *GitlabReleases) (err error) {
	if gitlabReleases == nil {
		return TracedErrorf("gitlabReleases is nil")
	}

	g.gitlabReleases = gitlabReleases

	return nil
}

func (g *GitlabRelease) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}
