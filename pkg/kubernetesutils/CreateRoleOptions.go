package kubernetesutils

import (
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CreateRoleOptions struct {
	Name     string
	Verbs    []string
	Resorces []string
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

func (c *CreateRoleOptions) GetVerbs() (verbs []string, err error) {
	if c.Verbs == nil {
		return nil, tracederrors.TracedErrorf("Verbs not set")
	}

	if len(c.Verbs) <= 0 {
		return nil, tracederrors.TracedErrorf("Verbs has no elements")
	}

	return c.Verbs, nil
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
