package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CreateClusterRoleOptions struct {
	Name      string
	Verbs     []string
	Resorces  []string
	APIGroups []string
}

func NewCreateClusterRoleOptions() (c *CreateClusterRoleOptions) {
	return new(CreateClusterRoleOptions)
}

func (c *CreateClusterRoleOptions) GetName() (name string, err error) {
	if c.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return c.Name, nil
}

func (c *CreateClusterRoleOptions) GetResorces() (resorces []string, err error) {
	if c.Resorces == nil {
		return nil, tracederrors.TracedErrorf("Resorces not set")
	}

	if len(c.Resorces) <= 0 {
		return nil, tracederrors.TracedErrorf("Resorces has no elements")
	}

	return c.Resorces, nil
}

func (c *CreateClusterRoleOptions) GetVerbs() (verbs []string, err error) {
	if c.Verbs == nil {
		return nil, tracederrors.TracedErrorf("Verbs not set")
	}

	if len(c.Verbs) <= 0 {
		return nil, tracederrors.TracedErrorf("Verbs has no elements")
	}

	return c.Verbs, nil
}

func (c *CreateClusterRoleOptions) GetAPIGroups() (apiGroups []string, err error) {
	if c.APIGroups == nil {
		return nil, tracederrors.TracedErrorf("APIGroups not set")
	}

	if len(c.APIGroups) <= 0 {
		return nil, tracederrors.TracedErrorf("APIGroups has no elements")
	}

	return c.APIGroups, nil
}

func (c *CreateClusterRoleOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.Name = name

	return nil
}

func (c *CreateClusterRoleOptions) SetResorces(resorces []string) (err error) {
	if resorces == nil {
		return tracederrors.TracedErrorf("resorces is nil")
	}

	if len(resorces) <= 0 {
		return tracederrors.TracedErrorf("resorces has no elements")
	}

	c.Resorces = resorces

	return nil
}

func (c *CreateClusterRoleOptions) SetVerbs(verbs []string) (err error) {
	if verbs == nil {
		return tracederrors.TracedErrorf("verbs is nil")
	}

	if len(verbs) <= 0 {
		return tracederrors.TracedErrorf("verbs has no elements")
	}

	c.Verbs = verbs

	return nil
}

func (c *CreateClusterRoleOptions) SetAPIGroups(apiGroups []string) (err error) {
	if apiGroups == nil {
		return tracederrors.TracedErrorf("apiGroups is nil")
	}

	if len(apiGroups) <= 0 {
		return tracederrors.TracedErrorf("apiGroups has no elements")
	}

	c.APIGroups = apiGroups

	return nil
}
