package yay

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func UpdateDatabase(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *packagemanageroptions.UpdateDatabaseOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	logging.LogInfoByCtxf(ctx, "Update yay database started.")

	if options == nil {
		options = new(packagemanageroptions.UpdateDatabaseOptions)
	}

	command := []string{"yay", "-Sy"}

	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	joinedCommand, err := shelllinehandler.Join(command)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Command to update yay database: '%s'.", joinedCommand)

	output, err := commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
		&parameteroptions.RunCommandOptions{
			Command:           command,
			AllowAllExitCodes: true,
		},
	)
	if err != nil {
		return err
	}
	if !output.IsExitSuccess() {
		stderr, err := output.GetStderrAsString()
		if err != nil {
			return err
		}

		return tracederrors.TracedErrorf("Failed to update yay database. The command '%s' failed with stderr: '%s'", joinedCommand, stderr)
	}

	logging.LogInfoByCtxf(ctx, "Update yay database finished.")
	return nil
}

func (p *Yay) UpdateDatabase(ctx context.Context, options *packagemanageroptions.UpdateDatabaseOptions) error {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	return UpdateDatabase(ctx, commandExecutor, options)
}
