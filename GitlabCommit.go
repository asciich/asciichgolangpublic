package asciichgolangpublic

type GitlabCommit struct {
	gitlabProjectsCommit *GitlabProjectCommits
}

func NewGitlabCommit() (g *GitlabCommit) {
	return new(GitlabCommit)
}

func (g *GitlabCommit) GetGitlabProjectsCommit() (gitlabProjectsCommit *GitlabProjectCommits, err error) {
	if g.gitlabProjectsCommit == nil {
		return nil, TracedErrorf("gitlabProjectsCommit not set")
	}

	return g.gitlabProjectsCommit, nil
}

func (g *GitlabCommit) MustGetGitlabProjectsCommit() (gitlabProjectsCommit *GitlabProjectCommits) {
	gitlabProjectsCommit, err := g.GetGitlabProjectsCommit()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProjectsCommit
}

func (g *GitlabCommit) MustSetGitlabProjectsCommit(gitlabProjectsCommit *GitlabProjectCommits) {
	err := g.SetGitlabProjectsCommit(gitlabProjectsCommit)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCommit) SetGitlabProjectsCommit(gitlabProjectsCommit *GitlabProjectCommits) (err error) {
	if gitlabProjectsCommit == nil {
		return TracedErrorf("gitlabProjectsCommit is nil")
	}

	g.gitlabProjectsCommit = gitlabProjectsCommit

	return nil
}
