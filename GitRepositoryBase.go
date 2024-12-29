package asciichgolangpublic

type GitRepositoryBase struct {
	parentRepositoryForBaseClass GitRepository
}

func NewGitRepositoryBase() (g *GitRepositoryBase) {
	return new(GitRepositoryBase)
}

func (g *GitRepositoryBase) CommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, TracedErrorNil("commitOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	createdCommit, err = parent.Commit(commitOptions)
	if err != nil {
		return nil, err
	}

	err = parent.Push(commitOptions.Verbose)
	if err != nil {
		return nil, err
	}

	return createdCommit, nil
}

func (g *GitRepositoryBase) CreateAndInit(createOptions *CreateRepositoryOptions) (err error) {
	if createOptions == nil {
		return TracedErrorNil("createOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	err = parent.Create(createOptions.Verbose)
	if err != nil {
		return err
	}

	err = parent.Init(createOptions)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitRepositoryBase) GetLatestTagVersion(verbose bool) (latestTagVersion Version, err error) {
	versionTags, err := g.ListVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	for _, tag := range versionTags {
		toCheck, err := tag.GetVersion()
		if err != nil {
			return nil, err
		}

		if latestTagVersion == nil {
			latestTagVersion = toCheck
		}

		latestTagVersion, err = Versions().ReturnNewerVersion(latestTagVersion, toCheck)
		if err != nil {
			return nil, err
		}
	}

	return latestTagVersion, nil
}

func (g *GitRepositoryBase) GetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository, err error) {
	if g.parentRepositoryForBaseClass == nil {
		return nil, TracedErrorf("parentRepositoryForBaseClass not set")
	}

	return g.parentRepositoryForBaseClass, nil
}

func (g *GitRepositoryBase) ListVersionTags(verbose bool) (versionTags []GitTag, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	allTags, err := parent.ListTags(verbose)
	if err != nil {
		return nil, err
	}

	versionTags = []GitTag{}
	for _, tag := range allTags {
		isVersionTag, err := tag.IsVersionTag()
		if err != nil {
			return nil, err
		}

		if isVersionTag {
			versionTags = append(versionTags, tag)
		}
	}

	return versionTags, nil
}

func (g *GitRepositoryBase) MustCommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := g.CommitAndPush(commitOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdCommit
}

func (g *GitRepositoryBase) MustCreateAndInit(createOptions *CreateRepositoryOptions) {
	err := g.CreateAndInit(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustGetLatestTagVersion(verbose bool) (latestTagVersion Version) {
	latestTagVersion, err := g.GetLatestTagVersion(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return latestTagVersion
}

func (g *GitRepositoryBase) MustGetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository) {
	parentRepositoryForBaseClass, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentRepositoryForBaseClass
}

func (g *GitRepositoryBase) MustListVersionTags(verbose bool) (versionTags []GitTag) {
	versionTags, err := g.ListVersionTags(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionTags
}

func (g *GitRepositoryBase) MustSetParentRepositoryForBaseClass(parentRepositoryForBaseClass GitRepository) {
	err := g.SetParentRepositoryForBaseClass(parentRepositoryForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) SetParentRepositoryForBaseClass(parentRepositoryForBaseClass GitRepository) (err error) {
	if parentRepositoryForBaseClass == nil {
		return TracedErrorf("parentRepositoryForBaseClass is nil")
	}

	g.parentRepositoryForBaseClass = parentRepositoryForBaseClass

	return nil
}
