package kubernetes

type KubernetesCluster interface {
	CreateNamespaceByName(namespaceName string, verbose bool) (createdNamespace Namespace, err error)
	DeleteNamespaceByName(namespaceName string, verbose bool) (err error)
	GetName() (name string, err error)
	GetNamespaceByName(name string) (namespace Namespace, err error)
	ListNamespaces(verbose bool) (namespaces []Namespace, err error)
	ListNamespaceNames(verbose bool) (namespaceNames []string, err error)
	MustNamespaceByNameExists(namespaceName string, verbose bool) (exist bool)
	MustDeleteNamespaceByName(namespaceName string, verbose bool)
	MustCreateNamespaceByName(namespaceName string, verbose bool) (createdNamespace Namespace)
	MustGetName() (name string)
	MustGetNamespaceByName(name string) (namespace Namespace)
	MustListNamespaces(verbose bool) (namespaces []Namespace)
	MustListNamespaceNames(verbose bool) (namespaceNames []string)
	NamespaceByNameExists(namespaceName string, verbose bool) (exist bool, err error)
}
