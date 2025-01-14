package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type GitlabCreateReleaseLinkOptions struct {
	Verbose bool
	Name    string
	Url     string
}

func NewGitlabCreateReleaseLinkOptions() (g *GitlabCreateReleaseLinkOptions) {
	return new(GitlabCreateReleaseLinkOptions)
}

func (g *GitlabCreateReleaseLinkOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", errors.TracedErrorf("Name not set")
	}

	return g.Name, nil
}

func (g *GitlabCreateReleaseLinkOptions) GetNameAndUrl() (name string, url string, err error) {
	name, err = g.GetName()
	if err != nil {
		return "", "", err
	}

	url, err = g.GetUrl()
	if err != nil {
		return "", "", err
	}

	return name, url, nil
}

func (g *GitlabCreateReleaseLinkOptions) GetUrl() (url string, err error) {
	if g.Url == "" {
		return "", errors.TracedErrorf("Url not set")
	}

	return g.Url, nil
}

func (g *GitlabCreateReleaseLinkOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabCreateReleaseLinkOptions) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabCreateReleaseLinkOptions) MustGetNameAndUrl() (name string, url string) {
	name, url, err := g.GetNameAndUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name, url
}

func (g *GitlabCreateReleaseLinkOptions) MustGetUrl() (url string) {
	url, err := g.GetUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return url
}

func (g *GitlabCreateReleaseLinkOptions) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateReleaseLinkOptions) MustSetUrl(url string) {
	err := g.SetUrl(url)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateReleaseLinkOptions) SetName(name string) (err error) {
	if name == "" {
		return errors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateReleaseLinkOptions) SetUrl(url string) (err error) {
	if url == "" {
		return errors.TracedErrorf("url is empty string")
	}

	g.Url = url

	return nil
}

func (g *GitlabCreateReleaseLinkOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
