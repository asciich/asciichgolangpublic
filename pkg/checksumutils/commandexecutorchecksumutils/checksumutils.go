package commandexecutorchecksumutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// getChecksumUsingCommand runs the given checksum command (e.g. "sha256sum")
// for the given path on the command executor's host and returns the parsed
// checksum (the first whitespace separated field of the output).
func getChecksumUsingCommand(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, checksumCommand string, path string) (checksum string, err error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	if checksumCommand == "" {
		return "", tracederrors.TracedErrorEmptyString("checksumCommand")
	}

	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Get %s of '%s' on '%s' started.", checksumCommand, path, hostDescription)

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{checksumCommand, path},
		},
	)
	if err != nil {
		return "", err
	}

	stdout = strings.TrimSpace(stdout)

	splitted := strings.Fields(stdout)
	if len(splitted) <= 0 {
		return "", tracederrors.TracedErrorf(
			"Unable to parse %s output for '%s' on '%s'. Output was empty.",
			checksumCommand,
			path,
			hostDescription,
		)
	}

	checksum = splitted[0]

	if checksum == "" {
		return "", tracederrors.TracedErrorf(
			"Parsed empty checksum from %s output for '%s' on '%s'.",
			checksumCommand,
			path,
			hostDescription,
		)
	}

	logging.LogInfoByCtxf(ctx, "Get %s of '%s' on '%s' finished. Checksum is '%s'.", checksumCommand, path, hostDescription, checksum)

	return checksum, nil
}

func GetMD5SumFromFileByPath(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (checksum string, err error) {
	return getChecksumUsingCommand(ctx, commandExecutor, "md5sum", path)
}

func GetSha1SumFromFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (checksum string, err error) {
	return getChecksumUsingCommand(ctx, commandExecutor, "sha1sum", path)
}

func GetSha256SumFromFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (checksum string, err error) {
	return getChecksumUsingCommand(ctx, commandExecutor, "sha256sum", path)
}

func GetSha512SumFromFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (checksum string, err error) {
	return getChecksumUsingCommand(ctx, commandExecutor, "sha512sum", path)
}
