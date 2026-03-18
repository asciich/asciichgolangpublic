package commandexecutorexec_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

// Run a simple command and check the output.
//
// Hint: For convenience there is a RunCommandAndGetStdoutAsString() function available
//
//	in case only stdout is needed after a successful exec.
func TestExecRunCommand(t *testing.T) {
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

				output, err := commandexecutorexec.RunCommand(
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

				output, err := commandexecutorexec.RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command:     tt.command,
						StdinString: tt.stdin,
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

// Test the convenience function RunCommandAndgetStdoutAsString to directly get the
// stdout after a successful exec.
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

				stdout, err := commandexecutorexec.RunCommandAndGetStdoutAsString(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedOutput, stdout)
			},
		)
	}
}

// Test the convenience function RunCommandAndgetStdoutAsString to directly get the
// stdout after a successful exec.
func TestExecRunCommandAndGetStdoutAsBytes(t *testing.T) {
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

				stdout, err := commandexecutorexec.RunCommandAndGetStdoutAsBytes(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, []byte(tt.expectedOutput), stdout)
			},
		)
	}
}

func Test_RunCommandAndGetStdoutAsIoReadCloser(t *testing.T) {
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

				readCloser, err := commandexecutorexec.RunCommandAndGetStdoutAsIoReadCloser(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)
				defer readCloser.Close()

				got, err := io.ReadAll(readCloser)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedOutput, string(got))
			},
		)
	}
}

func Test_RunCommandAndGetStdinAsIoWriteCloser(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"hello"},
		{"hello world\n"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempFile, err := tempfiles.CreateTemporaryFile(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

				writeCloser, err := commandexecutorexec.RunCommandAndGetStdinAsIoWriteCloser(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: []string{"tee", tempFile},
					},
				)
				require.NoError(t, err)
				defer writeCloser.Close()

				_, err = fmt.Fprint(writeCloser, tt.input)
				require.NoError(t, err)
				err = writeCloser.Close()
				require.NoError(t, err)

				got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
				require.NoError(t, err)
				require.EqualValues(t, tt.input, got)
			},
		)
	}
}
