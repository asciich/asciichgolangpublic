package yay

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RemovePackages(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageNames []string, options *packagemanageroptions.RemovePackageOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if len(packageNames) == 0 {
		return tracederrors.TracedError("packageNames is empty")
	}

	logging.LogInfoByCtxf(ctx, "Remove yay packages '%v' started.", packageNames)

	if options == nil {
		options = new(packagemanageroptions.RemovePackageOptions)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	isInstalled, err := IsPackagesInstalled(ctx, commandExecutor, packageNames)
	if err != nil {
		return err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Yay packages '%s' is installed. Going to remove it now.", packageNames)

		command := []string{"yay", "-R", "--noconfirm"}

		command = append(command, packageNames...)

		if options.UseSudo {
			command = append([]string{"sudo"}, command...)
		}

		stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
			commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfoByCtx(ctx, stdout)

		logging.LogChangedByCtxf(ctx, "Removed yay packages '%s'", packageNames)
	} else {
		logging.LogInfoByCtxf(ctx, "Yay packages '%s' is already absent on '%s'. Skip removal.", packageNames, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Remove yay package '%s' finished.", packageNames)

	return nil
}

func (p *Yay) RemovePackages(ctx context.Context, packageNames []string, options *packagemanageroptions.RemovePackageOptions) error {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	return RemovePackages(ctx, commandExecutor, packageNames, options)
}
