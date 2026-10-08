package commandexecutoruserutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsRunningAsRoot checks if the command executor is running as root user.
// It executes 'id -u' command and checks if the output is '0' (root user ID).
func IsRunningAsRoot(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"id", "-u"},
		},
	)
	if err != nil {
		return false, err
	}

	userId := strings.TrimSpace(output)
	isRunningAsRoot := userId == "0"

	if isRunningAsRoot {
		logging.LogInfoByCtxf(ctx, "Running as root since user id is '%s'.", userId)
	} else {
		logging.LogInfoByCtxf(ctx, "Not running as root, user id is '%s'.", userId)
	}

	return isRunningAsRoot, nil
}
