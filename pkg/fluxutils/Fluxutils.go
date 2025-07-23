package fluxutils

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/commandexecutor"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/fluxutils/commandexecutorflux"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/fluxutils/fluxinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/fluxutils/fluxparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/fluxutils/nativeflux"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func GetDefaultFluxImplementation() fluxinterfaces.Flux {
	return commandexecutorflux.NewcommandExecutorFlux(commandexecutor.Bash())
}

func InstallFlux(ctx context.Context, options *fluxparameteroptions.InstalFluxOptions) (fluxinterfaces.FluxDeployment, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	return GetDefaultFluxImplementation().InstallFlux(ctx, options)
}

func GetFluxDeployment(cluster kubernetesinterfaces.KubernetesCluster, namespaceName string) (fluxinterfaces.FluxDeployment, error) {
	if cluster == nil {
		return nil, tracederrors.TracedErrorNil("cluster")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	return nativeflux.GetFluxDeployment(cluster, namespaceName)
}