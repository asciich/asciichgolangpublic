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

func UpdatePackages(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageNames []string, options *packagemanageroptions.UpdatePackageOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if len(packageNames) == 0 {
		return tracederrors.TracedError("packageNames is empty")
	}

	logging.LogInfoByCtxf(ctx, "Update yay packages '%s' started.", packageNames)

	if options == nil {
		options = new(packagemanageroptions.UpdatePackageOptions)
	}

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

	isInstalled, err := IsPackagesInstalled(ctx, commandExecutor, packageNames)
	if err != nil {
		return err
	}

	if !isInstalled {
		err := InstallPackages(
			ctx,
			commandExecutor,
			packageNames,
			&packagemanageroptions.InstallPackageOptions{
				Force: options.Force,
			},
		)
		if err != nil {
			return err
		}
	}

	isUpdateAvailable, err := IsPackagesUpdateAvailable(
		ctx,
		commandExecutor,
		packageNames,
		&packagemanageroptions.UpdateDatabaseOptions{
			UseSudo: options.UseSudo,
		})
	if err != nil {
		return err
	}

	if isUpdateAvailable {
		logging.LogInfoByCtxf(ctx, "Going to update yay packages '%s'.", packageNames)

		command := []string{"yay", "-S", "--noconfirm"}

		if options.Force {
			command = append(command, "--overwrite='*'")
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

		logging.LogChangedByCtxf(ctx, "Yay package '%s' updated.", packageNames)
	} else {
		logging.LogInfoByCtxf(ctx, "Yay package '%s' is already up to date.", packageNames)
	}

	logging.LogInfoByCtxf(ctx, "Update yay package '%s' finished.", packageNames)

	return nil
}

func (p *Yay) UpdatePackages(ctx context.Context, packageNames []string, options *packagemanageroptions.UpdatePackageOptions) error {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	return UpdatePackages(ctx, commandExecutor, packageNames, options)
}
