package asciichgolangpublic

type GitlabDeleteReleaseOptions struct {
	DeleteCorrespondingTag bool
}

func NewGitlabDeleteReleaseOptions() (g *GitlabDeleteReleaseOptions) {
	return new(GitlabDeleteReleaseOptions)
}

func (g *GitlabDeleteReleaseOptions) GetDeleteCorrespondingTag() (deleteCorrespondingTag bool) {

	return g.DeleteCorrespondingTag
}

func (g *GitlabDeleteReleaseOptions) SetDeleteCorrespondingTag(deleteCorrespondingTag bool) {
	g.DeleteCorrespondingTag = deleteCorrespondingTag
}
