package helmutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils/helminterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils/helmparameteroptions"
)

func GetDefaultHelmImplementation() (helminterfaces.Helm, error) {
	return GetCommandExecutorHelm(commandexecutorbashoo.Bash())
}

func InstallHelmChart(ctx context.Context, options *helmparameteroptions.InstallHelmChartOptions) error {
	helm, err := GetDefaultHelmImplementation()
	if err != nil {
		return err
	}

	return helm.InstallHelmChart(ctx, options)
}
