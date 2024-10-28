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
				output2 := bash.MustRunCommandAndGetStdoutAsString(
					&RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
					},
				)

				assert.EqualValues(tt.expectedOutput, output)
				assert.EqualValues(tt.expectedOutput, output2)
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

				output2 := Bash().MustRunCommand(
					&RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
					},
				)

				assert.EqualValues(tt.expectedStdout, output2.MustGetStdoutAsString())
				assert.EqualValues(tt.expectedStderr, output2.MustGetStderrAsString())
				assert.EqualValues(0, output2.MustGetReturnCode())
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
				output2 := Bash().MustRunCommandAndGetStdoutAsFloat64(
					&RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
					},
				)

				assert.EqualValues(tt.expectedFloat64, output)
				assert.EqualValues(tt.expectedFloat64, output2)
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

				output := Bash().MustRunCommandAndGetStdoutAsLines(
					&RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)
				output2 := Bash().MustRunCommandAndGetStdoutAsLines(
					&RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
					},
				)

				assert.EqualValues(tt.expectedLines, output)
				assert.EqualValues(tt.expectedLines, output2)
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				output := Bash().MustRunOneLinerAndGetStdoutAsString(tt.oneLiner, verbose)
				assert.EqualValues(tt.expectedOutput, output)
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				output := Bash().MustRunCommand(&RunCommandOptions{
					Command: tt.command,
					Verbose: verbose,
				})

				assert.EqualValues(
					output.MustGetFirstLineOfStdoutAsString(),
					tt.expectedOutput,
				)
			},
		)
	}
}

func TestBashGetHostDescription(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					"localhost",
					Bash().MustGetHostDescription(),
				)
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				var bash CommandExecutor = Bash()
				output := bash.MustRunCommandAndGetStdoutAsBytes(
					&RunCommandOptions{
						Command:     tt.command,
						Verbose:     verbose,
						StdinString: tt.stdin,
					},
				)

				output2 := bash.MustRunCommandAndGetStdoutAsString(
					&RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
						StdinString:        tt.stdin,
					},
				)

				assert.EqualValues([]byte(tt.expectedOutput), output)
				assert.EqualValues(tt.expectedOutput, output2)
			},
		)
	}
}
