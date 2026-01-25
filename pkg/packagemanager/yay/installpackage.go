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

func InstallPackages(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageNames []string, options *packagemanageroptions.InstallPackageOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if len(packageNames) == 0 {
		return tracederrors.TracedError("packageNames is empty")
	}

	logging.LogInfoByCtxf(ctx, "Install yay packages '%v' started.", packageNames)

	if options == nil {
		options = new(packagemanageroptions.InstallPackageOptions)
	}

	isInstalled, err := IsPackagesInstalled(ctx, commandExecutor, packageNames)
	if err != nil {
		return err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Yay packages '%s' is already installed.", packageNames)
	} else {
		logging.LogInfoByCtxf(ctx, "Yay packages '%s' is not installed. Going to install it now.", packageNames)

		if options.UpdateDatabaseFirst {
			err := UpdateDatabase(
				ctx,
				commandExecutor,
				&packagemanageroptions.UpdateDatabaseOptions{
					UseSudo: options.UseSudo,
				},
			)
			if err != nil {
				return err
			}
		}

		command := []string{"yay", "-S", "--noconfirm"}

		if options.Force {
			command = append(command, "--overwrite=*")
		}

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

		logging.LogChangedByCtxf(ctx, "Installed yay packages '%s'", packageNames)
	}

	if options.UpdatePackage {
		err := UpdatePackages(
			ctx,
			commandExecutor,
			packageNames,
			&packagemanageroptions.UpdatePackageOptions{
				UpdateDatabaseFirst: false,
				Force:               options.Force,
				UseSudo:             options.UseSudo,
			},
		)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Install yay package '%s' finished.", packageNames)

	return nil
}

func (p *Yay) InstallPackages(ctx context.Context, packageNames []string, options *packagemanageroptions.InstallPackageOptions) error {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	return InstallPackages(ctx, commandExecutor, packageNames, options)
}
