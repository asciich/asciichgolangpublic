package kind

import "github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"

// Kubernetes in Docker
type Kind interface {
	ClusterByNameExists(clusterName string, verbose bool) (exists bool, err error)
	CreateClusterByName(clusterName string, verbose bool) (cluster kubernetesutils.KubernetesCluster, err error)
	DeleteClusterByName(clusterName string, verbose bool) (err error)
	GetClusterByName(clusterName string) (cluster kubernetesutils.KubernetesCluster, err error)
	ListClusterNames(verbose bool) (clusterNames []string, err error)
	MustClusterByNameExists(clusterName string, verbose bool) (exists bool)
	MustCreateClusterByName(clusterName string, verbose bool) (cluster kubernetesutils.KubernetesCluster)
	MustDeleteClusterByName(clusterName string, verbose bool)
	MustGetClusterByName(clusterName string) (cluster kubernetesutils.KubernetesCluster)
	MustListClusterNames(verbose bool) (clusterNames []string)
}
