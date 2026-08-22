package commandexecutorfile

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Move moves a file from src to dest using the command executor.
func Move(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, src string, dest string, options *filesoptions.MoveOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if src == "" {
		return tracederrors.TracedErrorEmptyString("src")
	}

	if dest == "" {
		return tracederrors.TracedErrorEmptyString("dest")
	}

	if options == nil {
		options = &filesoptions.MoveOptions{}
	}

	command := []string{"mv", src, dest}
	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Move '%s' to '%s' on '%s' started.", src, dest, hostDescription)

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to move '%s' to '%s' on '%s': %w", src, dest, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Moved '%s' to '%s' on '%s'.", src, dest, hostDescription)

	return nil
}

// Copy copies a file from src to dest using the command executor.
func Copy(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, src string, dest string, options *filesoptions.CopyOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if src == "" {
		return tracederrors.TracedErrorEmptyString("src")
	}

	if dest == "" {
		return tracederrors.TracedErrorEmptyString("dest")
	}

	if options == nil {
		options = &filesoptions.CopyOptions{}
	}

	command := []string{"cp", src, dest}
	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Copy '%s' to '%s' on '%s' started.", src, dest, hostDescription)

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to copy '%s' to '%s' on '%s': %w", src, dest, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Copied '%s' to '%s' on '%s'.", src, dest, hostDescription)

	return nil
}
