package kubernetes

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CreateRoleOptions struct {
	Name     string
	Verbs    []string
	Resorces []string
	Verbose  bool
}

func NewCreateRoleOptions() (c *CreateRoleOptions) {
	return new(CreateRoleOptions)
}

func (c *CreateRoleOptions) GetName() (name string, err error) {
	if c.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return c.Name, nil
}

func (c *CreateRoleOptions) GetResorces() (resorces []string, err error) {
	if c.Resorces == nil {
		return nil, tracederrors.TracedErrorf("Resorces not set")
	}

	if len(c.Resorces) <= 0 {
		return nil, tracederrors.TracedErrorf("Resorces has no elements")
	}

	return c.Resorces, nil
}

func (c *CreateRoleOptions) GetVerbose() (verbose bool) {

	return c.Verbose
}

func (c *CreateRoleOptions) GetVerbs() (verbs []string, err error) {
	if c.Verbs == nil {
		return nil, tracederrors.TracedErrorf("Verbs not set")
	}

	if len(c.Verbs) <= 0 {
		return nil, tracederrors.TracedErrorf("Verbs has no elements")
	}

	return c.Verbs, nil
}

func (c *CreateRoleOptions) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (c *CreateRoleOptions) MustGetResorces() (resorces []string) {
	resorces, err := c.GetResorces()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return resorces
}

func (c *CreateRoleOptions) MustGetVerbs() (verbs []string) {
	verbs, err := c.GetVerbs()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbs
}

func (c *CreateRoleOptions) MustSetName(name string) {
	err := c.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CreateRoleOptions) MustSetResorces(resorces []string) {
	err := c.SetResorces(resorces)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CreateRoleOptions) MustSetVerbs(verbs []string) {
	err := c.SetVerbs(verbs)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CreateRoleOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.Name = name

	return nil
}

func (c *CreateRoleOptions) SetResorces(resorces []string) (err error) {
	if resorces == nil {
		return tracederrors.TracedErrorf("resorces is nil")
	}

	if len(resorces) <= 0 {
		return tracederrors.TracedErrorf("resorces has no elements")
	}

	c.Resorces = resorces

	return nil
}

func (c *CreateRoleOptions) SetVerbose(verbose bool) {
	c.Verbose = verbose
}

func (c *CreateRoleOptions) SetVerbs(verbs []string) (err error) {
	if verbs == nil {
		return tracederrors.TracedErrorf("verbs is nil")
	}

	if len(verbs) <= 0 {
		return tracederrors.TracedErrorf("verbs has no elements")
	}

	c.Verbs = verbs

	return nil
}
