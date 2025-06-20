package nativeflux

import (
	"context"

	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type FluxDeployment struct {
	cluster *nativekubernetes.NativeKubernetesCluster

	namespace string
}

func GetFluxDeployment(cluster kubernetesinterfaces.KubernetesCluster, namespaceName string) (*FluxDeployment, error) {
	if cluster == nil {
		return nil, tracederrors.TracedErrorNil("cluster")
	}

	nativeCluster, ok := cluster.(*nativekubernetes.NativeKubernetesCluster)
	if !ok {
		return nil, tracederrors.TracedErrorf("cluster is not a native kubernetes cluster, it is of type: '%s'", datatypes.MustGetTypeName(cluster))
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	return &FluxDeployment{
		namespace: namespaceName,
		cluster:   nativeCluster,
	}, nil
}

func (f *FluxDeployment) GitRepositoryExists(ctx context.Context, name string, namespace string) (bool, error) {
	return false, tracederrors.TracedErrorNotImplemented()
}
