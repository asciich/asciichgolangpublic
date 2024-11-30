package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellLineHandlerSplit(t *testing.T) {
	tests := []struct {
		commandString    string
		expectedSplitted []string
	}{
		{"echo hello", []string{"echo", "hello"}},
		{"echo 'hello world'", []string{"echo", "hello world"}},
		{"echo \"hello world\"", []string{"echo", "hello world"}},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				splitted := ShellLineHandler().MustSplit(tt.commandString)
				assert.EqualValues(tt.expectedSplitted, splitted)
			},
		)
	}
}

func TestShellLineHandlerJoin(t *testing.T) {
	tests := []struct {
		command        []string
		expectedJoined string
	}{
		{[]string{"echo"}, "echo"},
		{[]string{"echo", ""}, "echo ''"},
		{[]string{"echo", " "}, "echo ' '"},
		{[]string{"echo", "abc\"abc"} , "echo 'abc\"abc'"}, // evalated using python -c "import shlex; print(shlex.join(['echo', 'abc\"abc']))"
		{[]string{"echo", "abc'abc"} , "echo 'abc'\"'\"'abc'"}, // evalated using python -c "import shlex; print(shlex.join(['echo', 'abc\'abc']))"
		{[]string{"echo", "hello"}, "echo hello"},
		{[]string{"echo", "hello world"}, "echo 'hello world'"},
		{[]string{"echo", "hello\nworld"}, "echo 'hello\nworld'"},
		{[]string{"echo", "hello\nworld\n"}, "echo 'hello\nworld\n'"},
		{[]string{"echo", "hello\\nworld\\n"}, "echo 'hello\\nworld\\n'"},
		{[]string{"echo", "hello \"world"}, "echo 'hello \"world'"},
		{[]string{"echo", "hello 'world"}, "echo 'hello '\"'\"'world'"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				joined := ShellLineHandler().MustJoin(tt.command)
				assert.EqualValues(tt.expectedJoined, joined)
			},
		)
	}
}

// TODO enable again func TestShellLineHandlerJoinTwice(t *testing.T) {
// TODO enable again 	tests := []struct {
// TODO enable again 		testcase string
// TODO enable again 	}{
// TODO enable again 		{"testcase"},
// TODO enable again 	}
// TODO enable again
// TODO enable again 	for _, tt := range tests {
// TODO enable again 		t.Run(
// TODO enable again 			MustFormatAsTestname(tt),
// TODO enable again 			func(t *testing.T) {
// TODO enable again 				assert := assert.New(t)
// TODO enable again
// TODO enable again 				joined1 := ShellLineHandler().MustJoin([]string{"echo", "hello \"world"})
// TODO enable again 				joined2 := ShellLineHandler().MustJoin([]string{"bash", "-c", joined1})
// TODO enable again
// TODO enable again 				expected := "bash -c 'echo '\"'\"'hello \"world'\"'\"''"
// TODO enable again 				assert.EqualValues(expected, joined2)
// TODO enable again
// TODO enable again 				for _, joined := range []string{joined1, joined2} {
// TODO enable again 					executedOutput := Shell().MustRunCommandAndGetStdoutAsString(&RunCommandOptions{Command: []string{"bash", "-c", joined}})
// TODO enable again 					executedOutput = strings.TrimSpace(executedOutput)
// TODO enable again 					assert.EqualValues("hello \"world", executedOutput)
// TODO enable again 				}
// TODO enable again 			},
// TODO enable again 		)
// TODO enable again 	}
// TODO enable again }

// TODO enable againfunc TestShellLineHandlerJoinThreeTimes(t *testing.T) {
// TODO enable again	tests := []struct {
// TODO enable again		testcase string
// TODO enable again	}{
// TODO enable again		{"testcase"},
// TODO enable again	}
// TODO enable again
// TODO enable again	for _, tt := range tests {
// TODO enable again		t.Run(
// TODO enable again			MustFormatAsTestname(tt),
// TODO enable again			func(t *testing.T) {
// TODO enable again				assert := assert.New(t)
// TODO enable again
// TODO enable again				joined1 := ShellLineHandler().MustJoin([]string{"echo", "hello \"world"})
// TODO enable again				joined2 := ShellLineHandler().MustJoin([]string{"bash", "-c", joined1})
// TODO enable again				joined3 := ShellLineHandler().MustJoin([]string{"bash", "-c", joined2})
// TODO enable again
// TODO enable again				expected := "bash -c 'bash -c '\"'\"'echo '\"'\"'\"'\"'\"'\"'\"'\"'hello \"world'\"'\"'\"'\"'\"'\"'\"'\"''\"'\"''"
// TODO enable again				assert.EqualValues(expected, joined3)
// TODO enable again
// TODO enable again				for _, joined := range []string{joined1, joined2, joined3} {
// TODO enable again					executedOutput := Shell().MustRunCommandAndGetStdoutAsString(&RunCommandOptions{Command: []string{"bash", "-c", joined}})
// TODO enable again					executedOutput = strings.TrimSpace(executedOutput)
// TODO enable again					assert.EqualValues("hello \"world", executedOutput)
// TODO enable again				}
// TODO enable again			},
// TODO enable again		)
// TODO enable again	}
// TODO enable again}
