package pacman

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (p *Pacman) UpdatePackage(ctx context.Context, packageName string) error {
	if packageName == "" {
		return tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Update pacman package '%s' started.", packageName)

	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"pacman", "-S", "--noconfirm", packageName},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Update pacman package '%s' finished.", packageName)

	return nil
}
