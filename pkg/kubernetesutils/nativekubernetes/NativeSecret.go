package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type NativeSecret struct {
	namespace *NativeNamespace
	name      string
}

func (n *NativeSecret) GetNamespace() (*NativeNamespace, error) {
	if n.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	return n.namespace, nil
}

func (n *NativeSecret) GetName() (string, error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *NativeSecret) Exists(ctx context.Context) (bool, error) {
	secretName, err := n.GetName()
	if err != nil {
		return false, nil
	}

	namespace, err := n.GetNamespace()
	if err != nil {
		return false, nil
	}

	return namespace.SecretByNameExists(ctx, secretName)
}
