package commandexecutorexecoo_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestExecRunCommandAndGetStdoutAsString(t *testing.T) {
	tests := []struct {
		command        []string
		expectedOutput string
	}{
		{[]string{"echo", "hello"}, "hello\n"},
		{[]string{"echo", "hello world"}, "hello world\n"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				var exec commandexecutorinterfaces.CommandExecutor = commandexecutorexecoo.Exec()
				output, err := exec.RunCommandAndGetStdoutAsString(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				output2, err := exec.RunCommandAndGetStdoutAsString(
					commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedOutput, output)
				require.EqualValues(t, tt.expectedOutput, output2)
			},
		)
	}
}

func TestExecRunCommandStdin(t *testing.T) {
	tests := []struct {
		stdin          string
		command        []string
		expectedOutput string
	}{
		{"abc", []string{"cat"}, "abc"},
		{"abc\n", []string{"cat"}, "abc\n"},
		{"abc \n", []string{"cat"}, "abc \n"},
		{"abc \n ", []string{"cat"}, "abc \n "},
		{" abc \n ", []string{"cat"}, " abc \n "},
		{"\n abc \n ", []string{"cat"}, "\n abc \n "},
		{"\n\n abc \n ", []string{"cat"}, "\n\n abc \n "},
		{"\n\n abc \n x", []string{"cat"}, "\n\n abc \n x"},
		{"x\n\n abc \n ", []string{"cat"}, "x\n\n abc \n "},
		{"\na\nb\nc\n", []string{"cat"}, "\na\nb\nc\n"},
		{"a\nb\nc\n", []string{"cat"}, "a\nb\nc\n"},
		{"a\nb\nc", []string{"cat"}, "a\nb\nc"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				var exec commandexecutorinterfaces.CommandExecutor = commandexecutorexecoo.Exec()
				output, err := exec.RunCommandAndGetStdoutAsBytes(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command:     tt.command,
						StdinString: tt.stdin,
					},
				)
				require.NoError(t, err)

				output2, err := exec.RunCommandAndGetStdoutAsString(
					commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
					&parameteroptions.RunCommandOptions{
						Command:     tt.command,
						StdinString: tt.stdin,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, []byte(tt.expectedOutput), output)
				require.EqualValues(t, tt.expectedOutput, output2)
			},
		)
	}
}

func TestExecEnvVar(t *testing.T) {
	t.Run("env var not set", func(t *testing.T) {
		ctx := getCtx()
		stdout, err := commandexecutorexecoo.Exec().RunCommandAndGetStdoutAsString(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "echo -en \"${MY_ENV}\""},
			},
		)
		require.NoError(t, err)
		require.Empty(t, stdout)
	})

	tests := []struct {
		value string
	}{
		{"a"},
		{"hello"},
		{"hello world"},
		{"HELLO WORLD"},
	}

	for _, tt := range tests {
		t.Run("env var set: "+tt.value, func(t *testing.T) {
			ctx := getCtx()
			stdout, err := commandexecutorexecoo.Exec().RunCommandAndGetStdoutAsString(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command: []string{"bash", "-c", "echo -en \"${MY_ENV}\""},
					AdditionalEnvVars: map[string]string{
						"MY_ENV": tt.value,
					},
				},
			)
			require.NoError(t, err)
			require.EqualValues(t, tt.value, stdout)

			// Addionally the PATH variable is check to ensure it's not overwritten or absent after defining AdditionalEnvVars:
			pathValue := os.Getenv("PATH")
			require.NotEmpty(t, pathValue)

			stdout, err = commandexecutorexecoo.Exec().RunCommandAndGetStdoutAsString(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command: []string{"bash", "-c", "echo -en \"${PATH}\""},
					AdditionalEnvVars: map[string]string{
						"MY_ENV": tt.value,
					},
				},
			)
			require.NoError(t, err)
			require.EqualValues(t, pathValue, stdout)
		})
	}
}
