package commandexecutorbashoo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
)

// This test case shows how to execute a simple command using exec.
func Test_EchoHelloWorld(t *testing.T) {
	// Get a context
	ctx := context.TODO()

	// Get the bash object
	bash := commandexecutorbashoo.Bash()

	// execute the command
	output, err := bash.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdout(ctx), // This optional step anables live output of stdout.
		&parameteroptions.RunCommandOptions{
			// This is the command to execute.
			// Every argument is appended as additional string to to the command slice:
			Command: []string{"echo", "hello", "world"},
		},
	)

	// We expect no error caused by our command:
	require.NoError(t, err)

	// Since there was no error the return code should be ok (=0)
	returnCode, err := output.GetReturnCode()
	require.NoError(t, err)
	require.EqualValues(t, 0, returnCode)
	// a shorter version would be:
	require.True(t, output.IsExitSuccess())

	// Stdout should contain "hello world"
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.EqualValues(t, "hello world\n", stdout)

	// stderr should be empty
	stderr, err := output.GetStderrAsString()
	require.NoError(t, err)
	require.EqualValues(t, "", stderr)
}
