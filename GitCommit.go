package asciichgolangpublic

type GitCommit struct {
	gitRepo GitRepository
	hash    string
}

func NewGitCommit() (g *GitCommit) {
	return new(GitCommit)
}

func (g *GitCommit) GetGitRepo() (gitRepo GitRepository, err error) {

	return g.gitRepo, nil
}

func (g *GitCommit) GetHash() (hash string, err error) {
	if g.hash == "" {
		return "", TracedErrorf("hash not set")
	}

	return g.hash, nil
}

func (g *GitCommit) MustGetGitRepo() (gitRepo GitRepository) {
	gitRepo, err := g.GetGitRepo()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitRepo
}

func (g *GitCommit) MustGetHash() (hash string) {
	hash, err := g.GetHash()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hash
}

func (g *GitCommit) MustSetGitRepo(gitRepo GitRepository) {
	err := g.SetGitRepo(gitRepo)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitCommit) MustSetHash(hash string) {
	err := g.SetHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitCommit) SetGitRepo(gitRepo GitRepository) (err error) {
	if gitRepo == nil {
		return TracedErrorNil("gitRepo")
	}

	g.gitRepo = gitRepo

	return nil
}

func (g *GitCommit) SetHash(hash string) (err error) {
	if hash == "" {
		return TracedErrorf("hash is empty string")
	}

	g.hash = hash

	return nil
}
