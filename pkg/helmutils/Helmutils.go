package helmutils

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/commandexecutor"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/helmutils/helminterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/helmutils/helmparameteroptions"
)

func GetDefaultHelmImplementation() (helminterfaces.Helm, error) {
	return GetCommandExecutorHelm(commandexecutor.Bash())
}

func InstallHelmChart(ctx context.Context, options *helmparameteroptions.InstallHelmChartOptions) error {
	helm, err := GetDefaultHelmImplementation()
	if err != nil {
		return err
	}

	return helm.InstallHelmChart(ctx, options)
}
