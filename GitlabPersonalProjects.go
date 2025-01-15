package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabPersonalProjects struct {
	gitlab *GitlabInstance
}

func NewGitlabPersonalProjects() (g *GitlabPersonalProjects) {
	return new(GitlabPersonalProjects)
}

func (g *GitlabPersonalProjects) GetGitlab() (gitlab *GitlabInstance, err error) {
	if g.gitlab == nil {
		return nil, tracederrors.TracedErrorf("gitlab not set")
	}

	return g.gitlab, nil
}

func (g *GitlabPersonalProjects) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabPersonalProjects) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabPersonalProjects) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedErrorf("gitlab is nil")
	}

	g.gitlab = gitlab

	return nil
}
