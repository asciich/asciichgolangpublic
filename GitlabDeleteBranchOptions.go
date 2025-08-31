package asciichgolangpublic

type GitlabDeleteBranchOptions struct {
	// By default the delete function waits until the deleted branch is not returned in the branch list anymore to avoid race conditions (branch is deleted but still listed by gitlab.)
	// SkipWaitForDeletion = true will skip this check/wait.
	SkipWaitForDeletion bool

}

func NewGitlabDeleteBranchOptions() (g *GitlabDeleteBranchOptions) {
	return new(GitlabDeleteBranchOptions)
}

func (g *GitlabDeleteBranchOptions) GetSkipWaitForDeletion() (skipWaitForDeletion bool) {

	return g.SkipWaitForDeletion
}

func (g *GitlabDeleteBranchOptions) SetSkipWaitForDeletion(skipWaitForDeletion bool) {
	g.SkipWaitForDeletion = skipWaitForDeletion
}
