package pacman

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (p *Pacman) UpdatePackage(ctx context.Context, packageName string, options *packagemanageroptions.UpdatePackageOptions) error {
	if packageName == "" {
		return tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Update pacman package '%s' started.", packageName)

	if options == nil {
		options = new(packagemanageroptions.UpdatePackageOptions)
	}

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

	isInstalled, err := p.IsPackageInstalled(ctx, packageName)
	if err != nil {
		return err
	}

	if !isInstalled {
		err := p.InstallPackage(
			ctx,
			packageName,
			&packagemanageroptions.InstallPackageOptions{
				Force: options.Force,
			},
		)
		if err != nil {
			return err
		}
	}

	isUpdateAvailable, err := p.IsPackageUpdateAvailalbe(
		ctx,
		packageName,
		&packagemanageroptions.UpdateDatabaseOptions{
			UseSudo: options.UseSudo,
		})
	if err != nil {
		return err
	}

	if isUpdateAvailable {
		logging.LogInfoByCtxf(ctx, "Going to update pacman package '%s'.", packageName)

		commandExecutor, err := p.GetCommandExecutor()
		if err != nil {
			return err
		}

		command := []string{"pacman", "-S", "--noconfirm"}

		if options.Force {
			command = append(command, "--overwrite='*'")
		}

		command = append(command, packageName)

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

		logging.LogChangedByCtxf(ctx, "Pacman package '%s' updated.", packageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Pacman package '%s' is already up to date.", packageName)
	}

	logging.LogInfoByCtxf(ctx, "Update pacman package '%s' finished.", packageName)

	return nil
}
