package kindutils

import "github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"

// Kubernetes in Docker
type Kind interface {
	ClusterByNameExists(clusterName string, verbose bool) (exists bool, err error)
	CreateClusterByName(clusterName string, verbose bool) (cluster kubernetesinterfaces.KubernetesCluster, err error)
	DeleteClusterByName(clusterName string, verbose bool) (err error)
	GetClusterByName(clusterName string) (cluster kubernetesinterfaces.KubernetesCluster, err error)
	ListClusterNames(verbose bool) (clusterNames []string, err error)
}
