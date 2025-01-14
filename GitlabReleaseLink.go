package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type GitlabReleaseLink struct {
	gitlabReleaseLinks *GitlabReleaseLinks
	name               string
	cachedUrl          string
}

func NewGitlabReleaseLink() (g *GitlabReleaseLink) {
	return new(GitlabReleaseLink)
}

func (g *GitlabReleaseLink) GetCachedUrl() (cachedUrl string, err error) {
	if g.cachedUrl == "" {
		return "", errors.TracedErrorf("cachedUrl not set")
	}

	return g.cachedUrl, nil
}

func (g *GitlabReleaseLink) GetGitlabReleaseLinks() (gitlabReleaseLinks *GitlabReleaseLinks, err error) {
	if g.gitlabReleaseLinks == nil {
		return nil, errors.TracedErrorf("gitlabReleaseLinks not set")
	}

	return g.gitlabReleaseLinks, nil
}

func (g *GitlabReleaseLink) GetName() (name string, err error) {
	if g.name == "" {
		return "", errors.TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GitlabReleaseLink) MustGetCachedUrl() (cachedUrl string) {
	cachedUrl, err := g.GetCachedUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cachedUrl
}

func (g *GitlabReleaseLink) MustGetGitlabReleaseLinks() (gitlabReleaseLinks *GitlabReleaseLinks) {
	gitlabReleaseLinks, err := g.GetGitlabReleaseLinks()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabReleaseLinks
}

func (g *GitlabReleaseLink) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabReleaseLink) MustSetCachedUrl(cachedUrl string) {
	err := g.SetCachedUrl(cachedUrl)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabReleaseLink) MustSetGitlabReleaseLinks(gitlabReleaseLinks *GitlabReleaseLinks) {
	err := g.SetGitlabReleaseLinks(gitlabReleaseLinks)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabReleaseLink) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabReleaseLink) SetCachedUrl(cachedUrl string) (err error) {
	if cachedUrl == "" {
		return errors.TracedErrorf("cachedUrl is empty string")
	}

	g.cachedUrl = cachedUrl

	return nil
}

func (g *GitlabReleaseLink) SetGitlabReleaseLinks(gitlabReleaseLinks *GitlabReleaseLinks) (err error) {
	if gitlabReleaseLinks == nil {
		return errors.TracedErrorf("gitlabReleaseLinks is nil")
	}

	g.gitlabReleaseLinks = gitlabReleaseLinks

	return nil
}

func (g *GitlabReleaseLink) SetName(name string) (err error) {
	if name == "" {
		return errors.TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}
