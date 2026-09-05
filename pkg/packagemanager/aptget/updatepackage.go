package aptget

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

	logging.LogInfoByCtxf(ctx, "Update apt-get packages '%s' started.", packageNames)

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
		logging.LogInfoByCtxf(ctx, "Going to update apt-get packages '%s'.", packageNames)

		// "--only-upgrade" ensures we only upgrade already installed packages and
		// do not accidentally install new ones here. DEBIAN_FRONTEND=noninteractive
		// avoids blocking on interactive prompts (e.g. tzdata):
		command := []string{"env", "DEBIAN_FRONTEND=noninteractive", "apt-get", "install", "-y", "--only-upgrade"}

		if options.Force {
			// Equivalent to pacman's "--overwrite='*'": allow dpkg to overwrite files
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
			commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfoByCtx(ctx, stdout)

		logging.LogChangedByCtxf(ctx, "Apt-get package '%s' updated.", packageNames)
	} else {
		logging.LogInfoByCtxf(ctx, "Apt-get package '%s' is already up to date.", packageNames)
	}

	logging.LogInfoByCtxf(ctx, "Update apt-get package '%s' finished.", packageNames)

	return nil
}

func (a *AptGet) UpdatePackages(ctx context.Context, packageNames []string, options *packagemanageroptions.UpdatePackageOptions) error {
	commandExecutor, err := a.GetCommandExecutor()
	if err != nil {
		return err
	}

	return UpdatePackages(ctx, commandExecutor, packageNames, options)
}
