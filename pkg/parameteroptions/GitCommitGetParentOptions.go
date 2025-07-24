package parameteroptions

type GitCommitGetParentsOptions struct {
	IncludeParentsOfParents bool
	Verbose                 bool
}

func NewGitCommitGetParentsOptions() (g *GitCommitGetParentsOptions) {
	return new(GitCommitGetParentsOptions)
}

func (g *GitCommitGetParentsOptions) GetIncludeParentsOfParents() (includeParentsOfParents bool) {

	return g.IncludeParentsOfParents
}

func (g *GitCommitGetParentsOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitCommitGetParentsOptions) SetIncludeParentsOfParents(includeParentsOfParents bool) {
	g.IncludeParentsOfParents = includeParentsOfParents
}

func (g *GitCommitGetParentsOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
