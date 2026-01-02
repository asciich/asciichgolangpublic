package pacman

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (p *Pacman) IsPackageUpdateAvailalbe(ctx context.Context, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	if packageName == "" {
		return false, tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Is pacman update available for package '%s' started.", packageName)

	if options == nil {
		options = new(packagemanageroptions.UpdateDatabaseOptions)
	}

	commandExectuor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	queryPackageUpdate := func() (*commandoutput.CommandOutput, error) {

		command := []string{"pacman", "-Qu", packageName}
		return commandExectuor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command:           command,
				AllowAllExitCodes: true,
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
			logging.LogInfoByCtxf(ctx, "Pacman database not present yet to check for updates for package '%s'. Going to download pacman database.", packageName)
			err := p.UpdateDatabase(
				ctx,
				options,
			)
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
			logging.LogInfoByCtxf(ctx, "Update for the package '%s' found using pacman.", packageName)
			return true, nil
		}
	}

	if isEmpty && exitCode == 1 {
		logging.LogInfoByCtxf(ctx, "No update for pacman package '%s' available.", packageName)
		return false, nil
	}

	stderr, err := output.GetStderrAsString()
	if err != nil {
		return false, err
	}

	return false, tracederrors.TracedErrorf("Failed to evaluate if an update of the pacman package '%s' is available: %s", packageName, stderr)
}
