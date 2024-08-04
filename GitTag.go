package asciichgolangpublic

type GitTag struct {
	gitRepository GitRepository
	hash          string
}

func NewGitTag() (g *GitTag) {
	return new(GitTag)
}

func (g *GitTag) GetGitRepository() (gitRepository GitRepository, err error) {

	return g.gitRepository, nil
}

func (g *GitTag) GetHash() (hash string, err error) {
	if g.hash == "" {
		return "", TracedErrorf("hash not set")
	}

	return g.hash, nil
}

func (g *GitTag) MustGetGitRepository() (gitRepository GitRepository) {
	gitRepository, err := g.GetGitRepository()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitRepository
}

func (g *GitTag) MustGetHash() (hash string) {
	hash, err := g.GetHash()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hash
}

func (g *GitTag) MustSetGitRepository(gitRepository GitRepository) {
	err := g.SetGitRepository(gitRepository)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitTag) MustSetHash(hash string) {
	err := g.SetHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitTag) SetGitRepository(gitRepository GitRepository) (err error) {
	g.gitRepository = gitRepository

	return nil
}

func (g *GitTag) SetHash(hash string) (err error) {
	if hash == "" {
		return TracedErrorf("hash is empty string")
	}

	g.hash = hash

	return nil
}
