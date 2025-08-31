package parameteroptions

type CreateRepositoryOptions struct {
	BareRepository            bool
	InitializeWithEmptyCommit bool

	// Set the default author for the repository to a default one.
	// Mainly usefull for testing since the author stays everywhere the same.
	InitializeWithDefaultAuthor bool
}

func NewCreateRepositoryOptions() (c *CreateRepositoryOptions) {
	return new(CreateRepositoryOptions)
}

func (c *CreateRepositoryOptions) GetBareRepository() (bareRepository bool) {

	return c.BareRepository
}

func (c *CreateRepositoryOptions) GetInitializeWithDefaultAuthor() (initializeWithDefaultAuthor bool) {

	return c.InitializeWithDefaultAuthor
}

func (c *CreateRepositoryOptions) GetInitializeWithEmptyCommit() (initializeWithEmptyCommit bool) {

	return c.InitializeWithEmptyCommit
}

func (c *CreateRepositoryOptions) SetBareRepository(bareRepository bool) {
	c.BareRepository = bareRepository
}

func (c *CreateRepositoryOptions) SetInitializeWithDefaultAuthor(initializeWithDefaultAuthor bool) {
	c.InitializeWithDefaultAuthor = initializeWithDefaultAuthor
}

func (c *CreateRepositoryOptions) SetInitializeWithEmptyCommit(initializeWithEmptyCommit bool) {
	c.InitializeWithEmptyCommit = initializeWithEmptyCommit
}
