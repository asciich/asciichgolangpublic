package asciichgolangpublic

type GitRepositoryBase struct {
	parentRepositoryForBaseClass GitRepository
}

func NewGitRepositoryBase() (g *GitRepositoryBase) {
	return new(GitRepositoryBase)
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

func (g *GitRepositoryBase) GetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository, err error) {
	if g.parentRepositoryForBaseClass == nil {
		return nil, TracedErrorf("parentRepositoryForBaseClass not set")
	}

	return g.parentRepositoryForBaseClass, nil
}

func (g *GitRepositoryBase) MustCreateAndInit(createOptions *CreateRepositoryOptions) {
	err := g.CreateAndInit(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustGetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository) {
	parentRepositoryForBaseClass, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentRepositoryForBaseClass
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
