package commandexecutorfile

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// AppendBytes appends bytes to a file using the command executor.
func AppendBytes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, toWrite []byte) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	// Use tee -a to append to file
	writer, err := commandExecutor.RunCommandAndGetStdinAsIoWriteCloser(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"tee", "-a", path},
		},
	)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to open writer for appending to '%s' on '%s': %w", path, hostDescription, err)
	}
	defer writer.Close()

	_, err = writer.Write(toWrite)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to append bytes to '%s' on '%s': %w", path, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Appended bytes to file '%s' on '%s'.", path, hostDescription)

	return nil
}

// AppendString appends a string to a file using the command executor.
func AppendString(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, toWrite string) error {
	return AppendBytes(ctx, commandExecutor, path, []byte(toWrite))
}
