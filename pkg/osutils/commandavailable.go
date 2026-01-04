package osutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Returns true if the command is available on the system the commandexecutor is running in:
func IsCommandAvailable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, command string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if command == "" {
		return false, tracederrors.TracedErrorEmptyString("command")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	cmd := []string{"whereis", command}
	cmdJoined, err := shelllinehandler.Join(cmd)
	if err != nil {
		return false, err
	}
	output, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return false, tracederrors.TracedErrorf("Failed to evaluate if command exists '%s': '%s' failed: %w", command, cmdJoined, err)
	}

	paths := strings.TrimPrefix(strings.TrimSpace(output), command+":")
	paths = strings.TrimSpace(paths)
	isAvailable := len(paths) > 0

	if isAvailable {
		logging.LogInfoByCtxf(ctx, "Command '%s' is available on '%s'.", command, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Command '%s' is not available on '%s'.", command, hostDescription)
	}

	return isAvailable, nil
}
