package pacman

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (p *Pacman) InstallPackage(ctx context.Context, packageName string, options *packagemanageroptions.InstallPackageOptions) error {
	if packageName == "" {
		return tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Install pacman package '%s' started.", packageName)

	if options == nil {
		options = new(packagemanageroptions.InstallPackageOptions)
	}

	isInstalled, err := p.IsPackageInstalled(ctx, packageName)
	if err != nil {
		return err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Pacman package '%s' is already installed.", packageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Pacman package '%s' is not installed. Going to install it now.", packageName)

		if options.UpdateDatabaseFirst {
			err := p.UpdateDatabase(
				ctx,
				&packagemanageroptions.UpdateDatabaseOptions{
					UseSudo: options.UseSudo,
				},
			)
			if err != nil {
				return err
			}
		}

		command := []string{"pacman", "-S", "--noconfirm"}

		if options.Force {
			command = append(command, "--overwrite=*")
		}

		command = append(command, packageName)

		if options.UseSudo {
			command = append([]string{"sudo"}, command...)
		}

		commandExecutor, err := p.GetCommandExecutor()
		if err != nil {
			return err
		}

		stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfoByCtx(ctx, stdout)

		logging.LogChangedByCtxf(ctx, "Installed pacman package '%s'", packageName)
	}

	if options.UpdatePackage {
		err := p.UpdatePackage(
			ctx,
			packageName,
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

	logging.LogInfoByCtxf(ctx, "Install pacman package '%s' finished.", packageName)

	return nil
}
