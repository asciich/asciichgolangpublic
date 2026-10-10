package commandexecutorfile

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// verifyingWriteCloser wraps the stdin io.WriteCloser of the 'tee' command.
//
// Not every CommandExecutor reports a non-zero exit code of the command when
// closing stdin (e.g. 'tee' failing because the file can not be created).
// Therefore Close() verifies the written file exists and its size matches the
// number of bytes written. The size is taken from the filesystem metadata
// (GetSizeBytes uses 'stat'), the file content is not read.
type verifyingWriteCloser struct {
	ctx             context.Context
	commandExecutor commandexecutorinterfaces.CommandExecutor
	path            string
	writeCloser     io.WriteCloser
	bytesWritten    int64
	closed          bool
}

func (v *verifyingWriteCloser) Write(p []byte) (int, error) {
	n, err := v.writeCloser.Write(p)
	v.bytesWritten += int64(n)
	return n, err
}

// Close closes the underlying writer and verifies the result.
// Calling Close multiple times is safe, only the first call has an effect.
func (v *verifyingWriteCloser) Close() error {
	if v.closed {
		return nil
	}
	v.closed = true

	err := v.writeCloser.Close()
	if err != nil {
		return tracederrors.TracedErrorf("Failed to close writer for '%s': %w", v.path, err)
	}

	sizeBytes, err := GetSizeBytes(v.ctx, v.commandExecutor, v.path)
	if err != nil {
		return tracederrors.TracedErrorf(
			"Failed to verify '%s' was written: %w",
			v.path,
			filesgeneric.GetAsError(err),
		)
	}

	if sizeBytes != v.bytesWritten {
		return tracederrors.TracedErrorf(
			"Writing '%s' failed: expected '%d' bytes but file has '%d' bytes",
			v.path,
			v.bytesWritten,
			sizeBytes,
		)
	}

	return nil
}

// OpenAsWriteCloser returns an io.WriteCloser writing to 'path'.
//
// Missing parent directories are created before the write command is started.
// This must happen BEFORE starting 'tee' since 'tee' opens (and creates) the file
// immediately on startup, not on the first write.
func OpenAsWriteCloser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, options *filesoptions.WriteOptions) (io.WriteCloser, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	err := CreateParentDirectory(ctx, commandExecutor, path, &filesoptions.CreateOptions{
		UseSudo: options.UseSudo,
	})
	if err != nil {
		return nil, err
	}

	command := []string{"tee", path}
	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	writeCloser, err := commandExecutor.RunCommandAndGetStdinAsIoWriteCloser(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return nil, err
	}

	return &verifyingWriteCloser{
		ctx:             ctx,
		commandExecutor: commandExecutor,
		path:            path,
		writeCloser:     writeCloser,
	}, nil
}

func WriteString(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, content string, options *filesoptions.WriteOptions) error {
	return WriteBytes(ctx, commandExecutor, path, []byte(content), options)
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

	// OpenAsWriteCloser also creates missing parent directories.
	writer, err := OpenAsWriteCloser(ctx, commandExecutor, path, options)
	if err != nil {
		return err
	}

	_, err = writer.Write(content)
	if err != nil {
		_ = writer.Close()
		return tracederrors.TracedErrorf("Failed to write bytes to '%s': %w", path, err)
	}

	// Close() also verifies the file was written correctly (see verifyingWriteCloser).
	return writer.Close()
}
