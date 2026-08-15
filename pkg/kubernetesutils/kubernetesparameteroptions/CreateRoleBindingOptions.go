package kubernetesparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CreateRoleBindingOptions struct {
	Name      string
	RoleRef   string
	Subjects  []string
	SubjectKind string // "User", "Group", or "ServiceAccount"
}

func NewCreateRoleBindingOptions() (c *CreateRoleBindingOptions) {
	return new(CreateRoleBindingOptions)
}

func (c *CreateRoleBindingOptions) GetName() (name string, err error) {
	if c.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return c.Name, nil
}

func (c *CreateRoleBindingOptions) GetRoleRef() (roleRef string, err error) {
	if c.RoleRef == "" {
		return "", tracederrors.TracedErrorf("RoleRef not set")
	}

	return c.RoleRef, nil
}

func (c *CreateRoleBindingOptions) GetSubjects() (subjects []string, err error) {
	if c.Subjects == nil {
		return nil, tracederrors.TracedErrorf("Subjects not set")
	}

	if len(c.Subjects) <= 0 {
		return nil, tracederrors.TracedErrorf("Subjects has no elements")
	}

	return c.Subjects, nil
}

func (c *CreateRoleBindingOptions) GetSubjectKind() (subjectKind string, err error) {
	if c.SubjectKind == "" {
		return "", tracederrors.TracedErrorf("SubjectKind not set")
	}

	return c.SubjectKind, nil
}

func (c *CreateRoleBindingOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.Name = name

	return nil
}

func (c *CreateRoleBindingOptions) SetRoleRef(roleRef string) (err error) {
	if roleRef == "" {
		return tracederrors.TracedErrorf("roleRef is empty string")
	}

	c.RoleRef = roleRef

	return nil
}

func (c *CreateRoleBindingOptions) SetSubjects(subjects []string) (err error) {
	if subjects == nil {
		return tracederrors.TracedErrorf("subjects is nil")
	}

	if len(subjects) <= 0 {
		return tracederrors.TracedErrorf("subjects has no elements")
	}

	c.Subjects = subjects

	return nil
}

func (c *CreateRoleBindingOptions) SetSubjectKind(subjectKind string) (err error) {
	if subjectKind == "" {
		return tracederrors.TracedErrorf("subjectKind is empty string")
	}

	c.SubjectKind = subjectKind

	return nil
}
