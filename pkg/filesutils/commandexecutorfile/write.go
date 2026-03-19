package commandexecutorfile

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func OpenAsWriteCloser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, options *filesoptions.WriteOptions) (io.WriteCloser, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorEmptyString("options")
	}

	command := []string{"tee", path}
	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	return commandExecutor.RunCommandAndGetStdinAsIoWriteCloser(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
}
