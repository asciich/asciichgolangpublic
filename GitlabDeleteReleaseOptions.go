package asciichgolangpublic

type GitlabDeleteReleaseOptions struct {
	Verbose                bool
	DeleteCorrespondingTag bool
}

func NewGitlabDeleteReleaseOptions() (g *GitlabDeleteReleaseOptions) {
	return new(GitlabDeleteReleaseOptions)
}

func (g *GitlabDeleteReleaseOptions) GetDeleteCorrespondingTag() (deleteCorrespondingTag bool) {

	return g.DeleteCorrespondingTag
}

func (g *GitlabDeleteReleaseOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabDeleteReleaseOptions) SetDeleteCorrespondingTag(deleteCorrespondingTag bool) {
	g.DeleteCorrespondingTag = deleteCorrespondingTag
}

func (g *GitlabDeleteReleaseOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
