package fluxinterfaces

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/fluxutils/fluxparameteroptions"
)

type Flux interface {
	InstallFlux(ctx context.Context, options *fluxparameteroptions.InstalFluxOptions) (FluxDeployment, error)
}
