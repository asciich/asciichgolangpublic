package helminterfaces

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/helmutils/helmparameteroptions"
)

type Helm interface {
	AddRepositoryByName(ctx context.Context, name string, url string) error
	InstallHelmChart(ctx context.Context, options *helmparameteroptions.InstallHelmChartOptions) error
}
