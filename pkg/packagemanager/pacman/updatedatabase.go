package pacman

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func (p *Pacman) UpdateDatabase(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Update pacman database started.")

	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"pacman", "-Sy"},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Update pacman database finished.")
	return nil
}
