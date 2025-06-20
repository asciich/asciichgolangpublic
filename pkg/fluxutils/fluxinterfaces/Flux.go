package fluxinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxparameteroptions"
)

type Flux interface {
	InstallFlux(ctx context.Context, options *fluxparameteroptions.InstalFluxOptions) (FluxDeployment, error)
}
