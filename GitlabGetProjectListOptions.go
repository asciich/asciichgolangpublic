package asciichgolangpublic

type GitlabgetProjectListOptions struct {
	Owned   bool
	Verbose bool
}

func NewGitlabgetProjectListOptions() (g *GitlabgetProjectListOptions) {
	return new(GitlabgetProjectListOptions)
}

func (g *GitlabgetProjectListOptions) GetOwned() (owned bool) {

	return g.Owned
}

func (g *GitlabgetProjectListOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabgetProjectListOptions) SetOwned(owned bool) {
	g.Owned = owned
}

func (g *GitlabgetProjectListOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
