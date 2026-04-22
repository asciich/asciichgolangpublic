package commandexecutorfile

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Truncate(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, newSizeBytes int64) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if newSizeBytes < 0 {
		return tracederrors.TracedErrorf(
			"Invalid size for truncating: newSizeBytes='%d'",
			newSizeBytes,
		)
	}

	currentSize, err := GetSizeBytes(ctx, commandExecutor, path)
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	if currentSize == newSizeBytes {
		logging.LogInfof(
			"File '%s' on host '%s' is already of size '%d' bytes. Skip truncate.",
			path,
			hostDescription,
			newSizeBytes,
		)
	} else {
		_, err = commandExecutor.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{
					"truncate",
					fmt.Sprintf("-s%d", newSizeBytes),
					path,
				},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"File '%s' on host '%s' is truncated to '%d' bytes.",
			path,
			hostDescription,
			newSizeBytes,
		)
	}

	return nil
}

func GetSizeBytes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (int64, error) {
	if commandExecutor == nil {
		return 0, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return 0, tracederrors.TracedErrorEmptyString("path")
	}

	fileSize, err := commandExecutor.RunCommandAndGetStdoutAsInt64(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"stat", "--printf=%s", path,
			},
		},
	)
	if err != nil {
		return -1, err
	}

	return fileSize, nil
}
