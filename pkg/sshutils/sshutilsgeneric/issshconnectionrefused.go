package sshutilsgeneric

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
)

// isSshConnectionRefused reports whether the given error / command output
// indicates that the underlying ssh connection was refused.
func IsSshConnectionRefused(err error, commandOutput *commandoutput.CommandOutput) bool {
	// Already tagged by a lower layer.
	if netutilserrors.IsConnectionRefusedError(err) {
		return true
	}

	if containsConnectionRefused(err.Error()) {
		return true
	}

	if commandOutput != nil {
		// Also inspect stderr of the ssh invocation, where the message actually appears.
		if commandOutput != nil {
			if stderr, stderrErr := commandOutput.GetStderrAsString(); stderrErr == nil {
				if containsConnectionRefused(stderr) {
					return true
				}
			}
		}
	}

	return false
}

func containsConnectionRefused(text string) bool {
	return strings.Contains(strings.ToLower(text), "connection refused")
}
