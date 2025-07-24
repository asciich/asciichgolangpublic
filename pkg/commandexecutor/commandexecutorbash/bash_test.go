package commandexecutorbash_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbash"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestBashRunCommandAndGetStdoutAsString(t *testing.T) {
	tests := []struct {
		command        []string
		expectedOutput string
	}{
		{[]string{"echo", "hello"}, "hello\n"},
		{[]string{"echo", "hello world"}, "hello world\n"},
		{[]string{"echo hello world"}, "hello world\n"},
		{[]string{"echo 'hello world'"}, "hello world\n"},
		{[]string{"true && echo yes || echo no"}, "yes\n"},
		{[]string{"false && echo yes || echo no"}, "no\n"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				output, err := commandexecutorbash.RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				stdout, err := output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedOutput, stdout)
			},
		)
	}
}

func TestBashRunCommand(t *testing.T) {
	tests := []struct {
		command        []string
		expectedStdout string
		expectedStderr string
	}{
		{[]string{"echo", "hello"}, "hello\n", ""},
		{[]string{"echo", "hello world"}, "hello world\n", ""},
		{[]string{"bash", "-c", "echo \"hello world\""}, "hello world\n", ""},
		{[]string{"bash", "-c", "echo 'hello world'"}, "hello world\n", ""},
		{[]string{"bash", "-c", "echo 'hello world' 1>&2"}, "", "hello world\n"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				output, err := commandexecutorbash.RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				stdout, err := output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedStdout, stdout)

				stderr, err := output.GetStderrAsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedStderr, stderr)

				returnCode, err := output.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, 0, returnCode)

				output2, err := commandexecutorbash.RunCommand(
					commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				stdout2, err := output2.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedStdout, stdout2)

				stderr2, err := output2.GetStderrAsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedStderr, stderr2)

				returnCode2, err := output2.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, 0, returnCode2)
			},
		)
	}
}

func TestBashRunCommandExitCode(t *testing.T) {
	type TestCase struct {
		command          []string
		expectedExitCode int
	}

	tests := []TestCase{}
	for i := 0; i < 128; i++ {
		tests = append(tests, TestCase{
			command:          []string{"bash", "-c", fmt.Sprintf("exit %v", i)},
			expectedExitCode: i,
		})
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				output, err := commandexecutorbash.RunCommand(
					getCtx(),
					&parameteroptions.RunCommandOptions{
						Command:           tt.command,
						AllowAllExitCodes: true,
					},
				)
				require.NoError(t, err)

				returnCode, err := output.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedExitCode, returnCode)
			},
		)
	}
}

func TestBashRunOneLinerAndGetStdoutAsString(t *testing.T) {
	tests := []struct {
		oneLiner       string
		expectedOutput string
	}{
		{"echo hallo", "hallo\n"},
		{"echo", "\n"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				output, err := commandexecutorbash.RunOneLinerAndGetStdoutAsString(getCtx(), tt.oneLiner)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedOutput, output)
			},
		)
	}
}

func TestBashCommandAndGetFirstLineOfStdoutAsString(t *testing.T) {
	tests := []struct {
		command        []string
		expectedOutput string
	}{
		{[]string{"echo hallo"}, "hallo"},
		{[]string{"echo -ne hallo"}, "hallo"},
		{[]string{"echo -ne 'hallo\ndu'"}, "hallo"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				output, err := commandexecutorbash.RunCommand(
					getCtx(),
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				firstLine, err := output.GetFirstLineOfStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, firstLine, tt.expectedOutput)
			},
		)
	}
}

func TestBashRunCommandStdin(t *testing.T) {
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				output, err := commandexecutorbash.RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command:     tt.command,
						StdinString: tt.stdin,
					},
				)
				require.NoError(t, err)

				stdout, err := output.GetStdoutAsBytes()
				require.NoError(t, err)
				require.EqualValues(t, []byte(tt.expectedOutput), stdout)
			},
		)
	}
}
