package fluxparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type InstalFluxOptions struct {
	KubernetesCluster kubernetesinterfaces.KubernetesCluster
	Namespace         string
}

func (i *InstalFluxOptions) GetKubernetesCluster() (kubernetesinterfaces.KubernetesCluster, error) {
	if i.KubernetesCluster == nil {
		return nil, tracederrors.TracedError("KubernetesCluster not set")
	}

	return i.KubernetesCluster, nil
}

func (i *InstalFluxOptions) GetNamespace() (string, error) {
	if i.Namespace == "" {
		return "", tracederrors.TracedError("Namespace not set")
	}

	return i.Namespace, nil
}
