package kubernetesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
)

type KubernetesCluster interface {
	CreateNamespaceByName(ctx context.Context, namespaceName string) (createdNamespace Namespace, err error)
	CreateSecret(ctx context.Context, namespaceName string, secretName string, options *CreateSecretOptions) (createdSecret Secret, err error)
	DeleteNamespaceByName(ctx context.Context, namespaceName string) (err error)
	DeleteSecretByName(ctx context.Context, namespaceName string, secretName string) (err error)
	GetKubectlContext(ctx context.Context) (contextName string, err error)
	GetName() (name string, err error)
	GetNamespaceByName(name string) (namespace Namespace, err error)
	GetResourceByNames(resourceName string, resourceType string, namespaceName string) (resource Resource, err error)
	ListNamespaces(ctx context.Context) (namespaces []Namespace, err error)
	ListNamespaceNames(ctx context.Context) (namespaceNames []string, err error)
	ListResources(options *parameteroptions.ListKubernetesResourcesOptions) (resources []Resource, err error)
	ListResourceNames(options *parameteroptions.ListKubernetesResourcesOptions) (resourceNames []string, err error)
	NamespaceByNameExists(ctx context.Context, namespaceName string) (exists bool, err error)
	SecretByNameExists(ctx context.Context, namespaceName string, secretName string) (exists bool, err error)
}
