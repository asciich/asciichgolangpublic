package parameteroptions

type GitCommitGetParentsOptions struct {
	IncludeParentsOfParents bool
}

func NewGitCommitGetParentsOptions() (g *GitCommitGetParentsOptions) {
	return new(GitCommitGetParentsOptions)
}

func (g *GitCommitGetParentsOptions) GetIncludeParentsOfParents() (includeParentsOfParents bool) {

	return g.IncludeParentsOfParents
}

func (g *GitCommitGetParentsOptions) SetIncludeParentsOfParents(includeParentsOfParents bool) {
	g.IncludeParentsOfParents = includeParentsOfParents
}
