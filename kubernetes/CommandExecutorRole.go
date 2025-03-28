package kubernetes

import (
	"github.com/asciich/asciichgolangpublic/logging"
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

func (c *CommandExecutorRole) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (c *CommandExecutorRole) MustGetNamespace() (namespace Namespace) {
	namespace, err := c.GetNamespace()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return namespace
}

func (c *CommandExecutorRole) MustSetName(name string) {
	err := c.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorRole) MustSetNamespace(namespace Namespace) {
	err := c.SetNamespace(namespace)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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
