package aptget

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
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

	logging.LogInfoByCtxf(ctx, "Install apt-get packages '%v' started.", packageNames)

	if options == nil {
		options = new(packagemanageroptions.InstallPackageOptions)
	}

	isInstalled, err := IsPackagesInstalled(ctx, commandExecutor, packageNames)
	if err != nil {
		return err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Apt-get packages '%s' is already installed.", packageNames)
	} else {
		logging.LogInfoByCtxf(ctx, "Apt-get packages '%s' is not installed. Going to install it now.", packageNames)

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

		// DEBIAN_FRONTEND=noninteractive avoids interactive prompts (e.g. tzdata)
		// which would otherwise block the non-interactive install:
		command := []string{"env", "DEBIAN_FRONTEND=noninteractive", "apt-get", "install", "-y"}

		if options.Force {
			// Equivalent to pacman's "--overwrite=*": allow dpkg to overwrite files
			// owned by other packages and permit downgrades/held package changes:
			command = append(command,
				"-o", "Dpkg::Options::=--force-overwrite",
				"--allow-downgrades",
				"--allow-change-held-packages",
			)
		}

		command = append(command, packageNames...)

		if options.UseSudo {
			command = append([]string{"sudo"}, command...)
		}

		stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
			contextutils.WithSilent(ctx),
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfoByCtx(ctx, stdout)

		logging.LogChangedByCtxf(ctx, "Installed apt-get packages '%s'", packageNames)
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

	logging.LogInfoByCtxf(ctx, "Install apt-get package '%s' finished.", packageNames)

	return nil
}

func (a *AptGet) InstallPackages(ctx context.Context, packageNames []string, options *packagemanageroptions.InstallPackageOptions) error {
	commandExecutor, err := a.GetCommandExecutor()
	if err != nil {
		return err
	}

	return InstallPackages(ctx, commandExecutor, packageNames, options)
}
