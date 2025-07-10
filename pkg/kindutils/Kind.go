package kindutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils/kindparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Kubernetes in Docker
type Kind interface {
	ClusterByNameExists(ctx context.Context, clusterName string) (exists bool, err error)
	CreateCluster(ctx context.Context, options *kindparameteroptions.CreateClusterOptions) (cluster kubernetesinterfaces.KubernetesCluster, err error)
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

	return kind.CreateCluster(ctx, &kindparameteroptions.CreateClusterOptions{Name: clusterName})
}

func DeleteClusterByName(ctx context.Context, clusterName string) error {
	if clusterName == "" {
		return tracederrors.TracedErrorEmptyString("clusterName")
	}

	kind, err := GetLocalCommandExecutorKind()
	if err != nil {
		return err
	}

	return kind.DeleteClusterByName(ctx, clusterName)
}

func DeleteClusterByNameIfInContinuousIntegration(ctx context.Context, clusterName string) error {
	if clusterName == "" {
		return tracederrors.TracedErrorEmptyString("clusterName")
	}

	if continuousintegration.IsRunningInContinuousIntegration() {
		return DeleteClusterByName(ctx, clusterName)
	}

	logging.LogInfoByCtxf(ctx, "Skip deletion of cluster '%s' since not running in continuous integration.", clusterName)

	return nil
}
