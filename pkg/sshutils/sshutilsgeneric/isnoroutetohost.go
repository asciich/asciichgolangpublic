package sshutilsgeneric

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
)

// IsNoRouteToHost reports whether the given error / command output
// indicates that the underlying ssh connection failed because there
// was no route to the target host.
func IsNoRouteToHost(err error, commandOutput *commandoutput.CommandOutput) bool {
	// Already tagged by a lower layer.
	if netutilserrors.IsNoRouteToHostError(err) {
		return true
	}

	if err != nil && containsNoRouteToHost(err.Error()) {
		return true
	}

	if commandOutput != nil {
		// Also inspect stderr of the ssh invocation, where the message actually appears.
		if stderr, stderrErr := commandOutput.GetStderrAsString(); stderrErr == nil {
			if containsNoRouteToHost(stderr) {
				return true
			}
		}
	}

	return false
}

func containsNoRouteToHost(text string) bool {
	return strings.Contains(strings.ToLower(text), "no route to host")
}
