package yay

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func IsPackagesUpdateAvailable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageNames []string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	if len(packageNames) == 0 {
		return false, tracederrors.TracedError("packageNames has no entries.")
	}

	logging.LogInfoByCtxf(ctx, "Is packaman packages update available for '%s' started.", packageNames)
	var updateAvailable bool
	var err error

	for _, name := range packageNames {
		updateAvailable, err = IsPackageUpdateAvailable(ctx, commandExecutor, name, options)
		if err != nil {
			return false, err
		}

		if updateAvailable {
			updateAvailable = true
			break
		}
	}

	logging.LogInfoByCtxf(ctx, "Is packaman packages update available for '%s' finished.", packageNames)

	return updateAvailable, nil
}

func IsPackageUpdateAvailable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if packageName == "" {
		return false, tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Is yay update available for package '%s' started.", packageName)

	if options == nil {
		options = new(packagemanageroptions.UpdateDatabaseOptions)
	}

	err := CreateYayInstallationUser(ctx, commandExecutor, options.UseSudo)
	if err != nil {
		return false, err
	}

	defer func() {
		err := DeleteYayInstallationUser(ctx, commandExecutor, options.UseSudo)
		if err != nil {
			logging.LogGoError(err)
		}
	}()

	command := []string{"yay", "-Qu", packageName}

	commandJoined, err := shelllinehandler.Join(command)
	if err != nil {
		return false, err
	}

	queryPackageUpdate := func() (*commandoutput.CommandOutput, error) {
		return commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command:            command,
				AllowAllExitCodes:  true,
				RunAsUser:          YAY_INSTALLATION_USER,
				UseSudoToRunAsUser: options.UseSudo,
			},
		)
	}

	output, err := queryPackageUpdate()
	if err != nil {
		return false, err
	}

	isEmpty, err := output.IsStdoutAndStderrEmpty()
	if err != nil {
		return false, err
	}

	exitCode, err := output.GetReturnCode()
	if err != nil {
		return false, err
	}

	if !isEmpty && exitCode == 1 {
		stderr, err := output.GetStderrAsString()
		if err != nil {
			return false, err
		}

		if stringsutils.ContainsAllIgnoreCase(stderr, []string{"database file for", "does not exist", "-Sy"}) {
			logging.LogInfoByCtxf(ctx, "Yay database not present yet to check for updates for package '%s'. Going to download yay database.", packageName)
			err := UpdateDatabase(ctx, commandExecutor, options)
			if err != nil {
				return false, err
			}
		}

		output, err = queryPackageUpdate()
		if err != nil {
			return false, err
		}

		isEmpty, err = output.IsStdoutAndStderrEmpty()
		if err != nil {
			return false, err
		}

		exitCode, err = output.GetReturnCode()
		if err != nil {
			return false, err
		}
	}

	if exitCode == 0 {
		stdout, err := output.GetStdoutAsString()
		if err != nil {
			return false, err
		}

		if strings.Contains(stdout, packageName) {
			logging.LogInfoByCtxf(ctx, "Update for the package '%s' found using yay.", packageName)
			return true, nil
		}
	}

	if isEmpty && exitCode == 1 {
		logging.LogInfoByCtxf(ctx, "No update for yay package '%s' available.", packageName)
		return false, nil
	}

	stderr, err := output.GetStderrAsString()
	if err != nil {
		return false, err
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return false, err
	}

	return false, tracederrors.TracedErrorf("Failed to evaluate if an update of the yay package '%s' is available using command='%s': stdout='%s' stderr='%s'", packageName, commandJoined, stdout, stderr)
}

func (p *Yay) IsPackageUpdateAvailable(ctx context.Context, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	return IsPackageUpdateAvailable(ctx, commandExecutor, packageName, options)
}
