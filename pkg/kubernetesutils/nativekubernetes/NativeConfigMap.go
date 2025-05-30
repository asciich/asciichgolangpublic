package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type NativeConfigMap struct {
	namespace *NativeNamespace
	name      string
}

func (n *NativeConfigMap) GetNamespace() (*NativeNamespace, error) {
	if n.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	return n.namespace, nil
}

func (n *NativeConfigMap) GetName() (string, error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *NativeConfigMap) Exists(ctx context.Context) (bool, error) {
	configMapName, err := n.GetName()
	if err != nil {
		return false, nil
	}

	namespace, err := n.GetNamespace()
	if err != nil {
		return false, nil
	}

	return namespace.ConfigMapByNameExists(ctx, configMapName)
}
