package fluxparameteroptions

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
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
