package commandexecutor_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/testutils"
)

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

				var bash commandexecutor.CommandExecutor = commandexecutor.Bash()
				output, err := bash.RunCommandAndGetStdoutAsString(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				output2, err := bash.RunCommandAndGetStdoutAsString(
					commandexecutor.WithLiveOutputOnStdout(ctx),
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

				output, err := commandexecutor.Bash().RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedStdout, output.MustGetStdoutAsString())
				require.EqualValues(t, tt.expectedStderr, output.MustGetStderrAsString())
				require.EqualValues(t, 0, output.MustGetReturnCode())

				output2, err := commandexecutor.Bash().RunCommand(
					commandexecutor.WithLiveOutputOnStdout(ctx),
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)

				require.EqualValues(t, tt.expectedStdout, output2.MustGetStdoutAsString())
				require.EqualValues(t, tt.expectedStderr, output2.MustGetStderrAsString())
				require.EqualValues(t, 0, output2.MustGetReturnCode())
			},
		)
	}
}

func TestBashRunCommandAndGetStdoutAsFloat64(t *testing.T) {
	tests := []struct {
		command         []string
		expectedFloat64 float64
	}{
		{[]string{"echo", "0"}, 0},
		{[]string{"echo", "1"}, 1.0},
		{[]string{"echo", "1.1"}, 1.1},
		{[]string{"echo", "-11.1"}, -11.1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				output, err := commandexecutor.Bash().RunCommandAndGetStdoutAsFloat64(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				output2, err := commandexecutor.Bash().RunCommandAndGetStdoutAsFloat64(
					commandexecutor.WithLiveOutputOnStdout(ctx),
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedFloat64, output)
				require.EqualValues(t, tt.expectedFloat64, output2)
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
				output, err := commandexecutor.Bash().RunCommand(
					getCtx(),
					&parameteroptions.RunCommandOptions{
						Command:           tt.command,
						AllowAllExitCodes: true,
					},
				)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedExitCode, output.MustGetReturnCode())
			},
		)
	}
}

func TestBashRunCommandAndGetStdoutAsLines(t *testing.T) {
	tests := []struct {
		command       []string
		expectedLines []string
	}{
		{[]string{"echo", "0"}, []string{"0"}},
		{[]string{"echo", "hello world"}, []string{"hello world"}},
		{[]string{"echo -en 'hello\\nworld'"}, []string{"hello", "world"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				output, err := commandexecutor.Bash().RunCommandAndGetStdoutAsLines(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				output2, err := commandexecutor.Bash().RunCommandAndGetStdoutAsLines(
					commandexecutor.WithLiveOutputOnStdout(ctx),
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedLines, output)
				require.EqualValues(t, tt.expectedLines, output2)
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
				output, err := commandexecutor.Bash().RunOneLinerAndGetStdoutAsString(getCtx(), tt.oneLiner)
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
				output, err := commandexecutor.Bash().RunCommand(
					getCtx(),
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				require.EqualValues(
					t,
					output.MustGetFirstLineOfStdoutAsString(),
					tt.expectedOutput,
				)
			},
		)
	}
}

func TestBashGetHostDescription(t *testing.T) {
	description, err := commandexecutor.Bash().GetHostDescription()
	require.NoError(t, err)
	require.EqualValues(t, "localhost", description)
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

				var bash commandexecutor.CommandExecutor = commandexecutor.Bash()
				output, err := bash.RunCommandAndGetStdoutAsBytes(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command:     tt.command,
						StdinString: tt.stdin,
					},
				)
				require.NoError(t, err)

				output2, err := bash.RunCommandAndGetStdoutAsString(
					commandexecutor.WithLiveOutputOnStdout(ctx),
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
