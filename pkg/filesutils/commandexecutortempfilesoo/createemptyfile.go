package commandexecutortempfilesoo

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateEmptyTemporaryFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (filesinterfaces.File, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"mktemp"},
		},
	)

	if err != nil {
		return nil, err
	}

	path := strings.TrimSpace(stdout)

	logging.LogChangedByCtxf(ctx, "Created empty temporary file '%s' on '%s'.", path, hostDescription)

	return commandexecutorfileoo.New(commandExecutor, path)
}

func CreateLocalEmptyTemporaryFile(ctx context.Context) (filesinterfaces.File, error) {
	return CreateEmptyTemporaryFile(ctx, commandexecutorexecoo.Exec())
}
