package commandexecutorkubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorClusterRole struct {
	name              string
	kubernetesCluster kubernetesinterfaces.KubernetesCluster
}

func NewCommandExecutorClusterRole() (c *CommandExecutorClusterRole) {
	return new(CommandExecutorClusterRole)
}

func (c *CommandExecutorClusterRole) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorClusterRole) GetKubernetesCluster() (kubernetesCluster kubernetesinterfaces.KubernetesCluster, err error) {

	return c.kubernetesCluster, nil
}

func (c *CommandExecutorClusterRole) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorClusterRole) SetKubernetesCluster(kubernetesCluster kubernetesinterfaces.KubernetesCluster) (err error) {
	c.kubernetesCluster = kubernetesCluster

	return nil
}

func (c *CommandExecutorClusterRole) Exists(ctx context.Context) (exists bool, err error) {
	roleName, err := c.GetName()
	if err != nil {
		return false, err
	}

	cluster, err := c.GetKubernetesCluster()
	if err != nil {
		return false, err
	}

	return cluster.ClusterRoleByNameExists(ctx, roleName)
}

func (c *CommandExecutorClusterRole) Delete(ctx context.Context) (err error) {
	roleName, err := c.GetName()
	if err != nil {
		return err
	}

	cluster, err := c.GetKubernetesCluster()
	if err != nil {
		return err
	}

	return cluster.DeleteClusterRoleByName(ctx, roleName)
}
