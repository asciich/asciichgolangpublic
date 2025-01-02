package asciichgolangpublic

// TODO: This should become the generic "CreateBranchOptions" to use everywhere.
// It then should replace GitlabCreateBranchOptions.
type CreateBranchOptions struct {
	// Name of the branch to create:
	Name string

	// Enable verbose output:
	Verbose bool
}

func NewCreateBranchOptions() (c *CreateBranchOptions) {
	return new(CreateBranchOptions)
}

func (c *CreateBranchOptions) GetName() (name string, err error) {
	if c.Name == "" {
		return "", TracedErrorf("Name not set")
	}

	return c.Name, nil
}

func (c *CreateBranchOptions) GetVerbose() (verbose bool) {

	return c.Verbose
}

func (c *CreateBranchOptions) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (c *CreateBranchOptions) MustSetName(name string) {
	err := c.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CreateBranchOptions) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	c.Name = name

	return nil
}

func (c *CreateBranchOptions) SetVerbose(verbose bool) {
	c.Verbose = verbose
}
