package commandexecutortempfile

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateEmptyTemporaryDirectory(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"mktemp", "-d"},
		},
	)

	if err != nil {
		return "", err
	}

	path := strings.TrimSpace(stdout)

	logging.LogChangedByCtxf(ctx, "Created empty temporary directory '%s' on '%s'.", path, hostDescription)

	return path, nil
}
