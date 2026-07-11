package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"k8s.io/client-go/kubernetes"
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
		return false, err
	}

	namespace, err := n.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespace.SecretByNameExists(ctx, secretName)
}

func (n *NativeSecret) GetNamespaceName() (string, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (n *NativeSecret) GetClientSet() (*kubernetes.Clientset, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetClientSet()
}

func (n *NativeSecret) Read(ctx context.Context) (map[string][]byte, error) {
	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	secretName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ReadSecret(ctx, clientset, namespaceName, secretName)
}

func (n *NativeSecret) Delete(ctx context.Context) error {
	clientset, err := n.GetClientSet()
	if err != nil {
		return err
	}

	namespaceName, err := n.GetNamespaceName()
	if err != nil {
		return err
	}

	secretName, err := n.GetName()
	if err != nil {
		return err
	}

	return nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
}
