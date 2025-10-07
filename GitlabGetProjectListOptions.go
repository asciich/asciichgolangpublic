package asciichgolangpublic

type GitlabgetProjectListOptions struct {
	Owned bool
}

func NewGitlabgetProjectListOptions() (g *GitlabgetProjectListOptions) {
	return new(GitlabgetProjectListOptions)
}

func (g *GitlabgetProjectListOptions) GetOwned() bool {

	return g.Owned
}

func (g *GitlabgetProjectListOptions) SetOwned(owned bool) {
	g.Owned = owned
}
