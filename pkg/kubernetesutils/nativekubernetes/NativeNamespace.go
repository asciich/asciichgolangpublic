package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NativeNamespace struct {
	name              string
	kubernetesCluster *NativeKubernetesCluster
}

func (n *NativeNamespace) GetKubernetesCluster() (*NativeKubernetesCluster, error) {
	if n.kubernetesCluster == nil {
		return nil, tracederrors.TracedError("kubernetesCluster not set")
	}

	return n.kubernetesCluster, nil
}

func (n *NativeNamespace) GetClientSet() (*kubernetes.Clientset, error) {
	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return nil, err
	}

	return cluster.GetClientSet()
}

func (n *NativeNamespace) Create(ctx context.Context) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) CreateRole(ctx context.Context, createOptions *kubernetesutils.CreateRoleOptions) (createdRole kubernetesutils.Role, err error) {
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

func (n *NativeNamespace) GetResourceByNames(resourceName string, resourceType string) (resource kubernetesutils.Resource, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) GetRoleByName(name string) (role kubernetesutils.Role, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) ListRoleNames(ctx context.Context) (roleNames []string, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) RoleByNameExists(ctx context.Context, name string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) SecretByNameExists(ctx context.Context, secretName string) (bool, error) {
	if secretName == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return false, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return false, err
	}

	var exists bool
	_, err = clientset.CoreV1().Secrets(namespaceName).Get(ctx, secretName, metav1.GetOptions{})
	if err == nil {
		exists = true
	} else {
		if !errors.IsNotFound(err) {
			return false, tracederrors.TracedErrorf("failed to get secret '%s' in namespace '%s': %w", secretName, namespaceName, err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' exists.", secretName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' does not exist.", secretName, namespaceName)
	}

	return exists, nil
}

func (n *NativeNamespace) DeleteSecretByName(ctx context.Context, secretName string) (err error) {
	if secretName == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return err
	}

	exists, err := n.SecretByNameExists(ctx, secretName)
	if err != nil {
		return err
	}

	if exists {
		clientset, err := n.GetClientSet()
		if err != nil {
			return err
		}

		err = clientset.CoreV1().Secrets(namespaceName).Delete(ctx, secretName, metav1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete secret '%s' in namespace '%s'.", secretName, namespaceName)
		}

		logging.LogChangedByCtxf(ctx, "Secret '%s' in namespace '%s' deleted.", secretName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' does not exist. Skip delete.", secretName, namespaceName)
	}

	return nil
}

func (n *NativeNamespace) GetSecretByName(name string) (secret kubernetesutils.Secret, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("secret")
	}

	return &NativeSecret{
		namespace: n,
		name:      name,
	}, nil
}

func (n *NativeNamespace) CreateSecret(ctx context.Context, secretName string, options *kubernetesutils.CreateSecretOptions) (createdSecret kubernetesutils.Secret, err error) {
	if secretName == "" {
		return nil, tracederrors.TracedErrorEmptyString("secret")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	exists, err := n.SecretByNameExists(ctx, secretName)
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, tracederrors.TracedError("Update existing secret not implemented")
	} else {
		clientset, err := n.GetClientSet()
		if err != nil {
			return nil, err
		}

		secretData, err := options.GetSecretData()
		if err != nil {
			return nil, err
		}

		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:   secretName,
				Labels: map[string]string{},
			},
			Data: secretData,
			Type: v1.SecretTypeOpaque,
		}

		_, err = clientset.CoreV1().Secrets(namespaceName).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to create secret '%s' in namespace '%s': %w", secretName, namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Created secret '%s' in kubernetes namespace '%s'.", secretName, namespaceName)
	}

	return n.GetSecretByName(secretName)
}
