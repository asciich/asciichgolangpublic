package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type NativeClusterRoleBinding struct {
	name              string
	kubernetesCluster kubernetesinterfaces.KubernetesCluster
}

func (n *NativeClusterRoleBinding) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return n.name, nil
}

func (n *NativeClusterRoleBinding) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	n.name = name

	return nil
}

func (n *NativeClusterRoleBinding) Exists(ctx context.Context) (exists bool, err error) {
	bindingName, err := n.GetName()
	if err != nil {
		return false, err
	}

	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return false, err
	}

	return cluster.ClusterRoleBindingByNameExists(ctx, bindingName)
}

func (n *NativeClusterRoleBinding) Delete(ctx context.Context) (err error) {
	bindingName, err := n.GetName()
	if err != nil {
		return err
	}

	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return err
	}

	return cluster.DeleteClusterRoleBindingByName(ctx, bindingName)
}

func (n *NativeClusterRoleBinding) GetKubernetesCluster() (kubernetesCluster kubernetesinterfaces.KubernetesCluster, err error) {
	if n.kubernetesCluster == nil {
		return nil, tracederrors.TracedErrorf("kubernetesCluster not set")
	}

	return n.kubernetesCluster, nil
}

func (n *NativeClusterRoleBinding) SetKubernetesCluster(kubernetesCluster kubernetesinterfaces.KubernetesCluster) (err error) {
	if kubernetesCluster == nil {
		return tracederrors.TracedErrorf("kubernetesCluster is nil")
	}

	n.kubernetesCluster = kubernetesCluster

	return nil
}
