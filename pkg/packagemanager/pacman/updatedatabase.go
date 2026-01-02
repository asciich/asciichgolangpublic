package pacman

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func (p *Pacman) UpdateDatabase(ctx context.Context, options *packagemanageroptions.UpdateDatabaseOptions) error {
	logging.LogInfoByCtxf(ctx, "Update pacman database started.")

	if options == nil {
		options = new(packagemanageroptions.UpdateDatabaseOptions)
	}

	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	command := []string{"pacman", "-Sy"}

	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Update pacman database finished.")
	return nil
}
