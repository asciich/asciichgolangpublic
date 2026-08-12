package commandexecutorkubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorRoleBinding struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func NewCommandExecutorRoleBinding() (c *CommandExecutorRoleBinding) {
	return new(CommandExecutorRoleBinding)
}

func (c *CommandExecutorRoleBinding) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorRoleBinding) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorRoleBinding) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorRoleBinding) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorRoleBinding) Exists(ctx context.Context) (exists bool, err error) {
	bindingName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespaceObj, err := c.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespaceObj.RoleBindingByNameExists(ctx, bindingName)
}

func (c *CommandExecutorRoleBinding) Delete(ctx context.Context) (err error) {
	bindingName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceObj, err := c.GetNamespace()
	if err != nil {
		return err
	}

	return namespaceObj.DeleteRoleBindingByName(ctx, bindingName)
}
