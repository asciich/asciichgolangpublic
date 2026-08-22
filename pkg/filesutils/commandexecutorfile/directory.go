package commandexecutorfile

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func IsEmptyDirectory(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	command := []string{"ls", "-A", path}

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return false, tracederrors.TracedErrorf("Failed to list directory '%s' on '%s': %w", path, hostDescription, err)
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return false, err
	}

	isEmpty := strings.TrimSpace(stdout) == ""

	if isEmpty {
		logging.LogInfoByCtxf(ctx, "The directory '%s' on '%s' is empty.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "The directory '%s' on '%s' is not empty.", path, hostDescription)
	}

	return isEmpty, nil
}
