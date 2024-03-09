package asciichgolangpublic

type CreateRepositoryOptions struct {
	BareRepository bool
	Verbose        bool
}

func NewCreateRepositoryOptions() (c *CreateRepositoryOptions) {
	return new(CreateRepositoryOptions)
}

func (c *CreateRepositoryOptions) GetBareRepository() (bareRepository bool) {

	return c.BareRepository
}

func (c *CreateRepositoryOptions) GetVerbose() (verbose bool) {

	return c.Verbose
}

func (c *CreateRepositoryOptions) SetBareRepository(bareRepository bool) {
	c.BareRepository = bareRepository
}

func (c *CreateRepositoryOptions) SetVerbose(verbose bool) {
	c.Verbose = verbose
}
