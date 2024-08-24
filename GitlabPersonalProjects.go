package asciichgolangpublic

type GitlabPersonalProjects struct {
	gitlab *GitlabInstance
}

func NewGitlabPersonalProjects() (g *GitlabPersonalProjects) {
	return new(GitlabPersonalProjects)
}

func (g *GitlabPersonalProjects) GetGitlab() (gitlab *GitlabInstance, err error) {
	if g.gitlab == nil {
		return nil, TracedErrorf("gitlab not set")
	}

	return g.gitlab, nil
}

func (g *GitlabPersonalProjects) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabPersonalProjects) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabPersonalProjects) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedErrorf("gitlab is nil")
	}

	g.gitlab = gitlab

	return nil
}
