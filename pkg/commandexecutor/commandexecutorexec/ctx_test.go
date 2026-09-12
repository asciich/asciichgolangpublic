package commandexecutorexec_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// TestRunCommand_ContextTimeoutAbortsCommand verifies that RunCommand honors
// the context deadline: a "sleep 3" command must be aborted by a context with a
// 1 second timeout instead of running for the full 3 seconds.
func TestRunCommand_ContextTimeoutAbortsCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tStart := time.Now()
	_, err := commandexecutorexec.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sleep", "3"},
		},
	)
	elapsed := time.Since(tStart)

	// The command must have been aborted by the context.
	require.Error(t, err)

	// It must return close to the 1 second timeout, well before the 3 seconds
	// the "sleep" command would otherwise take.
	require.GreaterOrEqual(t, elapsed, 1*time.Second)
	require.Less(t, elapsed, 3*time.Second)
}
