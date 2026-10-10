package commandexecutorfile

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func IsEmptyFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	sizeBytes, err := GetSizeBytes(ctx, commandExecutor, path)
	if err != nil {
		return false, err
	}

	return sizeBytes == 0, nil
}

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

// GetSizeBytes returns the size of the file at 'path' in bytes.
//
// The size is taken from the filesystem metadata using 'stat', so the file
// itself is never read (no performance issues for big files).
//
// 'stat -c %s' is used instead of the GNU-only 'stat --printf=%s' since
// '-c' is supported by both GNU coreutils and BusyBox (e.g. Alpine Linux).
func GetSizeBytes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (int64, error) {
	if commandExecutor == nil {
		return 0, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return 0, tracederrors.TracedErrorEmptyString("path")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"stat", "-c", "%s", "--", path,
			},
		},
	)
	if err != nil {
		err = filesgeneric.GetAsError(err)
		return 0, err
	}

	// 'stat -c' adds a trailing newline to the output (unlike '--printf'):
	stdout = strings.TrimSpace(stdout)

	fileSize, err := strconv.ParseInt(stdout, 10, 64)
	if err != nil {
		return 0, tracederrors.TracedErrorf(
			"Unable to parse file size of '%s' from stat output '%s': %w",
			path,
			stdout,
			err,
		)
	}

	if fileSize < 0 {
		return 0, tracederrors.TracedErrorf(
			"Invalid negative file size '%d' for '%s'",
			fileSize,
			path,
		)
	}

	return fileSize, nil
}
