package commandexecutor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/testutils"
)

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
				require := require.New(t)

				const verbose bool = true

				var exec CommandExecutor = Exec()
				output := exec.MustRunCommandAndGetStdoutAsString(
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)

				output2 := exec.MustRunCommandAndGetStdoutAsString(
					&parameteroptions.RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
					},
				)

				require.EqualValues(tt.expectedOutput, output)
				require.EqualValues(tt.expectedOutput, output2)
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
				require := require.New(t)

				const verbose bool = true

				var exec CommandExecutor = Exec()
				output := exec.MustRunCommandAndGetStdoutAsBytes(
					&parameteroptions.RunCommandOptions{
						Command:     tt.command,
						Verbose:     verbose,
						StdinString: tt.stdin,
					},
				)

				output2 := exec.MustRunCommandAndGetStdoutAsString(
					&parameteroptions.RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
						StdinString:        tt.stdin,
					},
				)

				require.EqualValues([]byte(tt.expectedOutput), output)
				require.EqualValues(tt.expectedOutput, output2)
			},
		)
	}
}
