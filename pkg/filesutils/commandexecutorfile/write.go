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

func WriteBytes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, content []byte, options *filesoptions.WriteOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if content == nil {
		return tracederrors.TracedErrorNil("content")
	}

	if options == nil {
		options = &filesoptions.WriteOptions{}
	}

	writer, err := OpenAsWriteCloser(ctx, commandExecutor, path, options)
	if err != nil {
		return err
	}

	_, err = writer.Write(content)
	if err != nil {
		_ = writer.Close()
		return tracederrors.TracedErrorf("Failed to write bytes to '%s': %w", path, err)
	}

	err = writer.Close()
	if err != nil {
		return tracederrors.TracedErrorf("Failed to close writer for '%s': %w", path, err)
	}

	return nil
}
