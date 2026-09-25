package gitgeneric

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Branch struct {
	name string
}

func (b *Branch) GetName() (string, error) {
	if b.name == "" {
		return "", tracederrors.TracedError("Branch name not set")
	}

	return b.name, nil
}

func (b *Branch) SetName(name string) error {
	if strings.TrimSpace(name) == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	b.name = name

	return nil
}

func (b *Branch) GetDeepCopy() gitinterfaces.Branch {
	deepCopy := NewBranch()

	deepCopy.name = b.name

	return deepCopy
}

func NewBranch() *Branch {
	return new(Branch)
}
