package helmparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type InstallHelmChartOptions struct {
	KubernetesCluster kubernetesinterfaces.KubernetesCluster
	ChartReference    string
	ChartUri          string
	Namespace         string
}

func (i *InstallHelmChartOptions) GetKubernetesCluster() (kubernetesinterfaces.KubernetesCluster, error) {
	if i.KubernetesCluster == nil {
		return nil, tracederrors.TracedError("KubernetesCluster not set")
	}

	return i.KubernetesCluster, nil
}

func (i *InstallHelmChartOptions) GetChartReference() (string, error) {
	if i.ChartReference == "" {
		return "", tracederrors.TracedError("ChartReference not set")
	}

	return i.ChartReference, nil
}

func (i *InstallHelmChartOptions) GetChartUri() (string, error) {
	if i.ChartUri == "" {
		return "", tracederrors.TracedError("ChartUri not set")
	}

	return i.ChartUri, nil
}

func (i *InstallHelmChartOptions) GetNamespace() (string, error) {
	if i.Namespace == "" {
		return "", tracederrors.TracedError("Namespace not set")
	}

	return i.Namespace, nil
}
