package commandexecutorfile

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetBirthDate returns the birth (creation) time of the file or directory at
// the given path on the host targeted by the command executor.
//
// It uses 'stat -c %W' which prints the birth time in seconds since the epoch.
// A value of 0 means the filesystem does not record a birth time; in that case
// an error is returned (mirroring the native STATX_BTIME behaviour).
func GetBirthDate(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (time.Time, error) {
	if commandExecutor == nil {
		return time.Time{}, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return time.Time{}, tracederrors.TracedErrorEmptyString("path")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return time.Time{}, err
	}

	// Do not rely on the exit code alone (constitution): run stat inside a
	// 'sh -c' that prints a well defined marker so we can distinguish a
	// successful lookup, a missing file and an unexpected error.
	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				"stat -c %W '" + path + "' 2>/dev/null || echo missing",
			},
		},
	)
	if err != nil {
		return time.Time{}, tracederrors.TracedErrorf("Failed to stat '%s' on '%s' to get birth date: %w", path, hostDescription, err)
	}

	value := strings.TrimSpace(stdout)

	if value == "missing" || value == "" {
		return time.Time{}, tracederrors.TracedErrorf("File '%s' does not exist on '%s' to get birth date.", path, hostDescription)
	}

	seconds, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, tracederrors.TracedErrorf("Failed to parse birth date value '%s' of '%s' on '%s': %w", value, path, hostDescription, err)
	}

	// 'stat -c %W' returns 0 when the birth time is unknown / unsupported.
	if seconds <= 0 {
		return time.Time{}, tracederrors.TracedErrorf("Birth date is not available for '%s' on '%s'.", path, hostDescription)
	}

	birthDate := time.Unix(seconds, 0).UTC()

	logging.LogInfoByCtxf(ctx, "Birth date of '%s' on '%s' is '%v'.", path, hostDescription, birthDate)

	return birthDate, nil
}

// GetBirthDateRFC3339 returns the birth (creation) time of the file or
// directory at the given path formatted as an RFC3339 string.
func GetBirthDateRFC3339(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	birthDate, err := GetBirthDate(ctx, commandExecutor, path)
	if err != nil {
		return "", err
	}

	ret := birthDate.Format(time.RFC3339)

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Birth date of '%s' on '%s' in RFC3339 is '%s'.", path, hostDescription, ret)

	return ret, nil
}
