package kubernetesinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
)

type KubernetesCluster interface {
	CheckAccessible(ctx context.Context) error
	ConfigMapByNameExists(ctx context.Context, namespaceName string, configMapName string) (exists bool, err error)
	CreateConfigMap(ctx context.Context, namespaceName string, configMapName string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap ConfigMap, err error)
	CreateNamespaceByName(ctx context.Context, namespaceName string) (createdNamespace Namespace, err error)
	CreateSecret(ctx context.Context, namespaceName string, secretName string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret Secret, err error)
	DeleteNamespaceByName(ctx context.Context, namespaceName string) (err error)
	DeleteSecretByName(ctx context.Context, namespaceName string, secretName string) (err error)
	GetKubectlContext(ctx context.Context) (contextName string, err error)
	GetName() (name string, err error)
	GetNamespaceByName(name string) (namespace Namespace, err error)
	GetResourceByNames(resourceName string, kind string, namespaceName string) (resource Resource, err error)
	ListNamespaces(ctx context.Context) (namespaces []Namespace, err error)
	ListNamespaceNames(ctx context.Context) (namespaceNames []string, err error)
	ListResources(options *parameteroptions.ListKubernetesResourcesOptions) (resources []Resource, err error)
	ListResourceNames(options *parameteroptions.ListKubernetesResourcesOptions) (resourceNames []string, err error)
	NamespaceByNameExists(ctx context.Context, namespaceName string) (exists bool, err error)
	SecretByNameExists(ctx context.Context, namespaceName string, secretName string) (exists bool, err error)
	WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, namespaceName string, options *kubernetesparameteroptions.WaitForPodsOptions) error
	WhoAmI(ctx context.Context) (*kubernetesimplementationindependend.UserInfo, error)
}
