package commandexecutorfile

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsStaticallyLinkedBinary checks if the given path points to a statically linked binary.
// It uses the `file` command to determine this.
// Returns true if the file is a statically linked binary, false otherwise.
// Returns an error if the file does not exist or if the check fails.
func IsStaticallyLinkedBinary(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	// Use the `file` command to check if the binary is statically linked
	// Statically linked binaries show "statically linked" in the file output
	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{"file", filePath},
		},
	)
	if err != nil {
		return false, tracederrors.TracedErrorf("Failed to run 'file' command: %w", err)
	}

	// Check if the output contains "statically linked"
	// Example output: "/path/to/binary: ELF 64-bit LSB executable, x86-64, statically linked, ..."
	isStaticallyLinked := strings.Contains(stdout, "statically linked")

	return isStaticallyLinked, nil
}
