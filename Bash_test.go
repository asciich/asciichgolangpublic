package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				var bash CommandExecutor = Bash()
				output := bash.MustRunCommandAndGetStdoutAsString(
					&RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)

				assert.EqualValues(tt.expectedOutput, output)
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				output := Bash().MustRunCommand(
					&RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)

				assert.EqualValues(tt.expectedStdout, output.MustGetStdoutAsString())
				assert.EqualValues(tt.expectedStderr, output.MustGetStderrAsString())
				assert.EqualValues(0, output.MustGetReturnCode())
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				output := Bash().MustRunCommandAndGetStdoutAsFloat64(
					&RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)

				assert.EqualValues(tt.expectedFloat64, output)
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
			command: []string{"bash", "-c", fmt.Sprintf("exit %v", i)},
		})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				output := Bash().MustRunCommand(
					&RunCommandOptions{
						Command:           tt.command,
						AllowAllExitCodes: true,
					},
				)

				assert.EqualValues(tt.expectedExitCode, output.MustGetReturnCode())
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				output := Bash().MustRunCommandAngGetStdoutAsLines(
					&RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)

				assert.EqualValues(tt.expectedLines, output)
			},
		)
	}
}
