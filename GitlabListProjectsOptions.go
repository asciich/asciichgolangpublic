package asciichgolangpublic

type GitlabListProjectsOptions struct {
	Recursive bool
	Verbose   bool
}

func NewGitlabListProjectsOptions() (g *GitlabListProjectsOptions) {
	return new(GitlabListProjectsOptions)
}

func (g *GitlabListProjectsOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabListProjectsOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
