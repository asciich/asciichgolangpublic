package pacman

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (p *Pacman) IsPackageInstalled(ctx context.Context, packageName string) (bool, error) {
	if packageName == "" {
		return false, tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Check pacman package '%s' is installed started.", packageName)

	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	var isInstalled bool
	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"pacman", "-Qs", packageName},
		},
	)
	if err == nil {
		isInstalled = true
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Pacman package '%s' is already installed.", packageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Pacman package '%s' is not installed.", packageName)
	}

	logging.LogInfoByCtxf(ctx, "Check pacman package '%s' is installed finished.", packageName)

	return isInstalled, nil
}
