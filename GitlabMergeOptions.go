package asciichgolangpublic

type GitlabMergeOptions struct {
	Verbose bool
}

func NewGitlabMergeOptions() (g *GitlabMergeOptions) {
	return new(GitlabMergeOptions)
}

func (g *GitlabMergeOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabMergeOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
