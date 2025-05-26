package kubernetesutils

import (
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorRole struct {
	name      string
	namespace Namespace
}

func NewCommandExecutorRole() (c *CommandExecutorRole) {
	return new(CommandExecutorRole)
}

func (c *CommandExecutorRole) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorRole) GetNamespace() (namespace Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorRole) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorRole) SetNamespace(namespace Namespace) (err error) {
	c.namespace = namespace

	return nil
}
