package commandexecutorfile

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Chown changes the ownership of a file or directory using the command executor.
func Chown(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, options *parameteroptions.ChownOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	userAndGroup, err := options.GetUserAndOptionallyGroupForChownCommand()
	if err != nil {
		return err
	}

	command := []string{"chown", userAndGroup, path}
	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Chown '%s' to '%s' on '%s' started.", path, userAndGroup, hostDescription)

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to chown '%s' to '%s' on '%s': %w", path, userAndGroup, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Changed ownership of '%s' to '%s' on '%s'.", path, userAndGroup, hostDescription)

	return nil
}
