package aptget

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Returns true when all packages are installed.
func IsPackagesInstalled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageNames []string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if len(packageNames) == 0 {
		return false, tracederrors.TracedError("packageNames is empty")
	}

	logging.LogInfoByCtxf(ctx, "Check if apt-get packages '%v' are installed started.", packageNames)

	for _, name := range packageNames {
		isInstalled, err := IsPackageInstalled(ctx, commandExecutor, name)
		if err != nil {
			return false, err
		}

		if !isInstalled {
			return false, nil
		}
	}

	logging.LogInfoByCtxf(ctx, "Check if apt-get packages '%v' are installed finished.", packageNames)

	return true, nil
}

func IsPackageInstalled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageName string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if packageName == "" {
		return false, tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Check apt-get package '%s' is installed started.", packageName)

	// apt-get itself has no query subcommand, so we use dpkg-query to read the
	// package status. A package counts as installed only when its status line
	// contains "install ok installed" (this excludes states like
	// "deinstall ok config-files" where config files remain but the package is
	// effectively removed):
	var isInstalled bool
	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.WithSilent(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"dpkg-query", "-W", "-f=${Status}", packageName},
		},
	)
	if err == nil && strings.Contains(stdout, "install ok installed") {
		isInstalled = true
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Apt-get package '%s' is already installed.", packageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Apt-get package '%s' is not installed.", packageName)
	}

	logging.LogInfoByCtxf(ctx, "Check apt-get package '%s' is installed finished.", packageName)

	return isInstalled, nil
}

func (a *AptGet) IsPackageInstalled(ctx context.Context, packageName string) (bool, error) {
	commandExecutor, err := a.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	return IsPackageInstalled(ctx, commandExecutor, packageName)
}
