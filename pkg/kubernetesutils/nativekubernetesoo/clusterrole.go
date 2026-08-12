package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type NativeClusterRole struct {
	name              string
	kubernetesCluster kubernetesinterfaces.KubernetesCluster
}

func (n *NativeClusterRole) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return n.name, nil
}

func (n *NativeClusterRole) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	n.name = name

	return nil
}

func (n *NativeClusterRole) Exists(ctx context.Context) (exists bool, err error) {
	roleName, err := n.GetName()
	if err != nil {
		return false, err
	}

	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return false, err
	}

	return cluster.ClusterRoleByNameExists(ctx, roleName)
}

func (n *NativeClusterRole) Delete(ctx context.Context) (err error) {
	roleName, err := n.GetName()
	if err != nil {
		return err
	}

	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return err
	}

	return cluster.DeleteClusterRoleByName(ctx, roleName)
}

func (n *NativeClusterRole) GetKubernetesCluster() (kubernetesCluster kubernetesinterfaces.KubernetesCluster, err error) {
	if n.kubernetesCluster == nil {
		return nil, tracederrors.TracedErrorf("kubernetesCluster not set")
	}

	return n.kubernetesCluster, nil
}

func (n *NativeClusterRole) SetKubernetesCluster(kubernetesCluster kubernetesinterfaces.KubernetesCluster) (err error) {
	if kubernetesCluster == nil {
		return tracederrors.TracedErrorf("kubernetesCluster is nil")
	}

	n.kubernetesCluster = kubernetesCluster

	return nil
}
