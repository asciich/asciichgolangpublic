package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateReleaseLinkOptions struct {
	Name string
	Url  string
}

func NewGitlabCreateReleaseLinkOptions() (g *GitlabCreateReleaseLinkOptions) {
	return new(GitlabCreateReleaseLinkOptions)
}

func (g *GitlabCreateReleaseLinkOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
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
		return "", tracederrors.TracedErrorf("Url not set")
	}

	return g.Url, nil
}

func (g *GitlabCreateReleaseLinkOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateReleaseLinkOptions) SetUrl(url string) (err error) {
	if url == "" {
		return tracederrors.TracedErrorf("url is empty string")
	}

	g.Url = url

	return nil
}
