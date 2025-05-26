package kubernetesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type NativeNamespace struct {
	name              string
	kubernetesCluster *NativeKubernetesCluster
}

func (n *NativeNamespace) Create(ctx context.Context) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) CreateRole(ctx context.Context, createOptions *CreateRoleOptions) (createdRole Role, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) DeleteRoleByName(ctx context.Context, name string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) GetClusterName() (clusterName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) GetKubectlContext(ctx context.Context) (contextName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}
func (n *NativeNamespace) GetResourceByNames(resourceName string, resourceType string) (resource Resource, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) GetRoleByName(name string) (role Role, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) ListRoleNames(ctx context.Context) (roleNames []string, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
func (n *NativeNamespace) RoleByNameExists(ctx context.Context, name string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}
