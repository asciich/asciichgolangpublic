package ansibleutils

import (
	"errors"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

var ErrAnsibleGroupNotFound = errors.New("ansible group not found")

type AnsibleGroup struct {
	name string
}

func NewAnsibleGroupByName(name string) (g *AnsibleGroup, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	g = new(AnsibleGroup)

	err = g.SetGroupName(name)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (a *AnsibleGroup) Name() (name string) {
	return a.name
}

func (a *AnsibleGroup) SetGroupName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	a.name = name

	return nil
}

func (a *AnsibleGroup) GetGroupName() (name string, err error) {
	if a.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return a.name, nil
}
