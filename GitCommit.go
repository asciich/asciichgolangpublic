package asciichgolangpublic

type GitCommit struct {
	gitRepo GitRepository
	hash    string
}

func NewGitCommit() (g *GitCommit) {
	return new(GitCommit)
}

func (g *GitCommit) GetAgeSeconds() (age float64, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return -1, err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return -1, err
	}

	age, err = repo.GetCommitAgeSecondsByCommitHash(hash)
	if err != nil {
		return
	}

	return age, nil
}

func (g *GitCommit) GetAuthorEmail() (authorEmail string, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return "", err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return "", err
	}

	authorEmail, err = repo.GetAuthorEmailByCommitHash(hash)
	if err != nil {
		return
	}

	return authorEmail, nil
}

func (g *GitCommit) GetAuthorString() (authorString string, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return "", err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return "", err
	}

	authorString, err = repo.GetAuthorStringByCommitHash(hash)
	if err != nil {
		return
	}

	return authorString, nil
}

func (g *GitCommit) GetCommitMessage() (commitMessage string, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return "", err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return "", err
	}

	commitMessage, err = repo.GetCommitMessageByCommitHash(hash)
	if err != nil {
		return
	}

	return commitMessage, nil
}

func (g *GitCommit) GetGitRepo() (gitRepo GitRepository, err error) {
	if g.gitRepo == nil {
		return nil, TracedError("gitRepo not set")
	}
	return g.gitRepo, nil
}

func (g *GitCommit) GetHash() (hash string, err error) {
	if g.hash == "" {
		return "", TracedErrorf("hash not set")
	}

	return g.hash, nil
}

func (g *GitCommit) GetParentCommits(options *GitCommitGetParentsOptions) (parentCommit []*GitCommit, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return nil, err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, err
	}

	parentCommit, err = repo.GetCommitParentsByCommitHash(hash, options)
	if err != nil {
		return
	}

	return parentCommit, nil
}

func (g *GitCommit) HasParentCommit() (hasParentCommit bool, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return false, err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return false, err
	}

	hasParentCommit, err = repo.CommitHasParentCommitByCommitHash(hash)
	if err != nil {
		return false, err
	}

	return hasParentCommit, nil
}

func (g *GitCommit) MustGetAgeSeconds() (age float64) {
	age, err := g.GetAgeSeconds()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return age
}

func (g *GitCommit) MustGetAuthorEmail() (authorEmail string) {
	authorEmail, err := g.GetAuthorEmail()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authorEmail
}

func (g *GitCommit) MustGetAuthorString() (authorString string) {
	authorString, err := g.GetAuthorString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authorString
}

func (g *GitCommit) MustGetCommitMessage() (commitMessage string) {
	commitMessage, err := g.GetCommitMessage()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitMessage
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

func (g *GitCommit) MustGetParentCommits(options *GitCommitGetParentsOptions) (parentCommit []*GitCommit) {
	parentCommit, err := g.GetParentCommits(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentCommit
}

func (g *GitCommit) MustHasParentCommit() (hasParentCommit bool) {
	hasParentCommit, err := g.HasParentCommit()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasParentCommit
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
