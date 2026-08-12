package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type NativeRoleBinding struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func (n *NativeRoleBinding) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return n.name, nil
}

func (n *NativeRoleBinding) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {
	if n.namespace == nil {
		return nil, tracederrors.TracedErrorf("namespace not set")
	}

	return n.namespace, nil
}

func (n *NativeRoleBinding) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	n.name = name

	return nil
}

func (n *NativeRoleBinding) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	if namespace == nil {
		return tracederrors.TracedErrorf("namespace is nil")
	}

	n.namespace = namespace

	return nil
}

func (n *NativeRoleBinding) Exists(ctx context.Context) (exists bool, err error) {
	bindingName, err := n.GetName()
	if err != nil {
		return false, err
	}

	namespaceObj, err := n.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespaceObj.RoleBindingByNameExists(ctx, bindingName)
}

func (n *NativeRoleBinding) Delete(ctx context.Context) (err error) {
	bindingName, err := n.GetName()
	if err != nil {
		return err
	}

	namespaceObj, err := n.GetNamespace()
	if err != nil {
		return err
	}

	return namespaceObj.DeleteRoleBindingByName(ctx, bindingName)
}
