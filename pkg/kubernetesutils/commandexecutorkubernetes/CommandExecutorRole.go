package commandexecutorkubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorRole struct {
	name      string
	namespace kubernetesinterfaces.Namespace
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

func (c *CommandExecutorRole) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorRole) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorRole) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorRole) Exists(ctx context.Context) (exists bool, err error) {
	roleName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespaceObj, err := c.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespaceObj.RoleByNameExists(ctx, roleName)
}

func (c *CommandExecutorRole) Delete(ctx context.Context) (err error) {
	roleName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceObj, err := c.GetNamespace()
	if err != nil {
		return err
	}

	return namespaceObj.DeleteRoleByName(ctx, roleName)
}
