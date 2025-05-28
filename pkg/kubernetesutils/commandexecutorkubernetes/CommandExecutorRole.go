package commandexecutorkubernetes

import (
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorRole struct {
	name      string
	namespace kubernetesutils.Namespace
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

func (c *CommandExecutorRole) GetNamespace() (namespace kubernetesutils.Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorRole) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorRole) SetNamespace(namespace kubernetesutils.Namespace) (err error) {
	c.namespace = namespace

	return nil
}
