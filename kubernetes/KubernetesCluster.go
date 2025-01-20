package kubernetes

type KubernetesCluster interface {
	CreateNamespaceByName(namespaceName string, verbose bool) (createdNamespace Namespace, err error)
	DeleteNamespaceByName(namespaceName string, verbose bool) (err error)
	GetKubectlContext(verbose bool) (contextName string, err error)
	GetName() (name string, err error)
	GetNamespaceByName(name string) (namespace Namespace, err error)
	GetResourceByNames(resourceName string, resourceType string, namespaceName string) (resource Resource, err error)
	ListNamespaces(verbose bool) (namespaces []Namespace, err error)
	ListNamespaceNames(verbose bool) (namespaceNames []string, err error)
	MustNamespaceByNameExists(namespaceName string, verbose bool) (exist bool)
	MustDeleteNamespaceByName(namespaceName string, verbose bool)
	MustCreateNamespaceByName(namespaceName string, verbose bool) (createdNamespace Namespace)
	MustGetKubectlContext(verbose bool) (contextName string)
	MustGetName() (name string)
	MustGetNamespaceByName(name string) (namespace Namespace)
	MustGetResourceByNames(resourceName string, resourceType string, namespaceName string) (resource Resource)
	MustListNamespaces(verbose bool) (namespaces []Namespace)
	MustListNamespaceNames(verbose bool) (namespaceNames []string)
	NamespaceByNameExists(namespaceName string, verbose bool) (exist bool, err error)
}
