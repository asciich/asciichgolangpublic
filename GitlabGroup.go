package asciichgolangpublic

import (
	"github.com/xanzy/go-gitlab"
)

type GitlabGroup struct {
	gitlab *GitlabInstance
	id     int
}

func NewGitlabGroup() (gitlabGroup *GitlabGroup) {
	return new(GitlabGroup)
}

func (g *GitlabGroup) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabGroup) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabGroup) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabGroup) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabGroup) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabGroup) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *GitlabGroup) GetFqdn() (fqdn string, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return "", err
	}

	fqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (p *GitlabGroup) GetGitlab() (gitlab *GitlabInstance, err error) {
	if p.gitlab == nil {
		return nil, TracedError("gitlab is not set")
	}

	return p.gitlab, nil
}

func (p *GitlabGroup) GetId() (id int, err error) {
	if p.id <= 0 {
		return -1, TracedError("id not set")
	}

	return p.id, nil
}

func (p *GitlabGroup) GetNativeClient() (nativeClient *gitlab.Client, err error) {
	gitlab, err := p.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (p *GitlabGroup) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedError("gitlab is nil")
	}

	p.gitlab = gitlab

	return nil
}

func (p *GitlabGroup) SetId(id int) (err error) {
	if id <= 0 {
		return TracedErrorf("invalid id = '%d'", id)
	}

	p.id = id

	return nil
}
