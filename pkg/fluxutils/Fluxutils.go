package fluxutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/commandexecutorflux"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/nativeflux"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
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
