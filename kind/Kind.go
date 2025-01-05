package kind

import "github.com/asciich/asciichgolangpublic/kubernetes"

// Kubernetes in Docker
type Kind interface {
	ClusterByNameExists(clusterName string, verbose bool) (exists bool, err error)
	CreateClusterByName(clusterName string, verbose bool) (cluster kubernetes.KubernetesCluster, err error)
	DeleteClusterByName(clusterName string, verbose bool) (err error)
	ListClusterNames(verbose bool) (clusterNames []string, err error)
	MustClusterByNameExists(clusterName string, verbose bool) (exists bool)
	MustCreateClusterByName(clusterName string, verbose bool) (cluster kubernetes.KubernetesCluster)
	MustDeleteClusterByName(clusterName string, verbose bool)
	MustListClusterNames(verbose bool) (clusterNames []string)
}
