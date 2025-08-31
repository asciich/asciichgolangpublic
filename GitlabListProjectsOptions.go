package asciichgolangpublic

type GitlabListProjectsOptions struct {
	Recursive bool
}

func NewGitlabListProjectsOptions() (g *GitlabListProjectsOptions) {
	return new(GitlabListProjectsOptions)
}

