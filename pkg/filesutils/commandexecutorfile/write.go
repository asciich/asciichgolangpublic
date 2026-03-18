package commandexecutorfile

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func OpenAsWriteCloser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (io.WriteCloser, error) {
	if commandExecutor==nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return commandExecutor.RunCommandAndGetStdinAsIoWriteCloser(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"tee", path},
		},
	)
}
