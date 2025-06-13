package fluxutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/commandexecutorflux"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxparameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func GetDefaultFluxImplementation() fluxinterfaces.Flux {
	return commandexecutorflux.NewcommandExecutorFlux(commandexecutor.Bash())
}

func InstallFlux(ctx context.Context, options *fluxparameteroptions.InstalFluxOptions) (fluxinterfaces.DeployedFlux, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	return GetDefaultFluxImplementation().InstallFlux(ctx, options)
}
