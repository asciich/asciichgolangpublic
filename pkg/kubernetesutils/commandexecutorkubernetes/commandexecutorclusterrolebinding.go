package commandexecutorkubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorClusterRoleBinding struct {
	name              string
	kubernetesCluster kubernetesinterfaces.KubernetesCluster
}

func NewCommandExecutorClusterRoleBinding() (c *CommandExecutorClusterRoleBinding) {
	return new(CommandExecutorClusterRoleBinding)
}

func (c *CommandExecutorClusterRoleBinding) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorClusterRoleBinding) GetKubernetesCluster() (kubernetesCluster kubernetesinterfaces.KubernetesCluster, err error) {

	return c.kubernetesCluster, nil
}

func (c *CommandExecutorClusterRoleBinding) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorClusterRoleBinding) SetKubernetesCluster(kubernetesCluster kubernetesinterfaces.KubernetesCluster) (err error) {
	c.kubernetesCluster = kubernetesCluster

	return nil
}

func (c *CommandExecutorClusterRoleBinding) Exists(ctx context.Context) (exists bool, err error) {
	bindingName, err := c.GetName()
	if err != nil {
		return false, err
	}

	cluster, err := c.GetKubernetesCluster()
	if err != nil {
		return false, err
	}

	return cluster.ClusterRoleBindingByNameExists(ctx, bindingName)
}

func (c *CommandExecutorClusterRoleBinding) Delete(ctx context.Context) (err error) {
	bindingName, err := c.GetName()
	if err != nil {
		return err
	}

	cluster, err := c.GetKubernetesCluster()
	if err != nil {
		return err
	}

	return cluster.DeleteClusterRoleBindingByName(ctx, bindingName)
}
