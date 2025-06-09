package kindutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Kubernetes in Docker
type Kind interface {
	ClusterByNameExists(ctx context.Context, clusterName string) (exists bool, err error)
	CreateClusterByName(ctx context.Context, clusterName string) (cluster kubernetesinterfaces.KubernetesCluster, err error)
	DeleteClusterByName(ctx context.Context, clusterName string) (err error)
	GetClusterByName(clusterName string) (cluster kubernetesinterfaces.KubernetesCluster, err error)
	ListClusterNames(ctx context.Context) (clusterNames []string, err error)
}

func CreateCluster(ctx context.Context, clusterName string) (kubernetesinterfaces.KubernetesCluster, error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	kind, err := GetLocalCommandExecutorKind()
	if err != nil {
		return nil, err
	}

	return kind.CreateClusterByName(ctx, clusterName)
}
