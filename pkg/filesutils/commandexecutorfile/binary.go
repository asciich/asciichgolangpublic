package commandexecutorfile

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// staticallyLinkedIndicators are parts of the 'file' output indicating a statically linked binary.
//
//   - "statically linked": classic static binaries.
//     Example: "ELF 64-bit LSB executable, x86-64, ..., statically linked, stripped"
//   - "static-pie linked": static position independent executables (e.g. '/bin/busybox.static'
//     on alpine). These do NOT contain "statically linked" in the 'file' output.
//     Example: "ELF 64-bit LSB pie executable, x86-64, ..., static-pie linked, stripped"
var staticallyLinkedIndicators = []string{
	"statically linked",
	"static-pie linked",
}

// isStaticallyLinkedFileOutput evaluates the output of 'file -b' for a binary.
func isStaticallyLinkedFileOutput(fileOutput string) bool {
	for _, indicator := range staticallyLinkedIndicators {
		if strings.Contains(fileOutput, indicator) {
			return true
		}
	}

	return false
}

// IsStaticallyLinkedBinary checks if the given path points to a statically linked binary.
// It uses the `file` command to determine this.
// Returns true if the file is a statically linked binary, false otherwise.
// Returns an error if the file does not exist or if the check fails.
//
// The 'file' command is called with:
//   - '-L': Follow symlinks. Otherwise e.g. '/bin/ls' on alpine (a symlink to '/bin/busybox')
//     is only reported as "symbolic link to /bin/busybox" and the real binary is never evaluated.
//   - '-b': Omit the file name in the output so the path itself can not cause false matches
//     (e.g. a path like '/tmp/statically linked/foo').
//
// A nonexistent file returns an error detectable by filesgeneric.IsErrFileNotFound.
func IsStaticallyLinkedBinary(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{"file", "-L", "-b", "--", filePath},
		},
	)
	if err != nil {
		return false, tracederrors.TracedErrorf(
			"Failed to run 'file' command (is 'file' installed? e.g. 'apk add file' on alpine): %w",
			filesgeneric.GetAsError(err),
		)
	}

	output := strings.TrimSpace(stdout)

	// 'file' exits with 0 for nonexistent files but reports:
	//   "cannot open `/nonexistent/file' (No such file or directory)"
	if strings.HasPrefix(output, "cannot open") {
		return false, filesgeneric.GetAsError(
			tracederrors.TracedErrorf("Unable to evaluate if '%s' is statically linked: %s", filePath, output),
		)
	}

	isStaticallyLinked := isStaticallyLinkedFileOutput(output)

	if isStaticallyLinked {
		logging.LogInfoByCtxf(ctx, "'%s' is a statically linked binary: '%s'.", filePath, output)
	} else {
		logging.LogInfoByCtxf(ctx, "'%s' is not a statically linked binary: '%s'.", filePath, output)
	}

	return isStaticallyLinked, nil
}
