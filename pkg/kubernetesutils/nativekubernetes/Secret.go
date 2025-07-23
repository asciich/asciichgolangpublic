package nativekubernetes

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	return ReadSecret(ctx, clientset, namespaceName, secretName)
}

func ReadSecret(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, secretName string) (map[string][]byte, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if secretName == "" {
		return nil, tracederrors.TracedErrorEmptyString("secretName")
	}

	secret, err := clientset.CoreV1().Secrets(namespaceName).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, tracederrors.TracedErrorf("Unable to read secret '%s'. Secret does not exist in namespace '%s'", secretName, namespaceName)
		}
	}

	return secret.Data, nil
}
