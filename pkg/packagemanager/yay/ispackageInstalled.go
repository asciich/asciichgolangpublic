package yay

import (
	"context"

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

	logging.LogInfoByCtxf(ctx, "Check if yay packages '%v' are installed started.", packageNames)

	for _, name := range packageNames {
		isInstalled, err := IsPackageInstalled(ctx, commandExecutor, name)
		if err != nil {
			return false, err
		}

		if !isInstalled {
			return false, nil
		}
	}

	logging.LogInfoByCtxf(ctx, "Check if yay packages '%v' are installed finished.", packageNames)

	return true, nil
}

func IsPackageInstalled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageName string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if packageName == "" {
		return false, tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Check yay package '%s' is installed started.", packageName)

	var isInstalled bool
	_, err := commandExecutor.RunCommand(
		contextutils.WithSilent(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"yay", "-Qs", "^" + packageName + "$"},
		},
	)
	if err == nil {
		isInstalled = true
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Yay package '%s' is already installed.", packageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Yay package '%s' is not installed.", packageName)
	}

	logging.LogInfoByCtxf(ctx, "Check yay package '%s' is installed finished.", packageName)

	return isInstalled, nil
}

func (p *Yay) IsPackageInstalled(ctx context.Context, packageName string) (bool, error) {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	return IsPackageInstalled(ctx, commandExecutor, packageName)
}
