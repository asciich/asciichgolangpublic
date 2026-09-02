package shelllinehandler_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				splitted, err := shelllinehandler.Split(tt.commandString)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedSplitted, splitted)
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
		{[]string{"echo", "abc\"abc"}, "echo 'abc\"abc'"},
		{[]string{"echo", "abc'abc"}, "echo 'abc'\"'\"'abc'"},
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				joined, err := shelllinehandler.Join(tt.command)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedJoined, joined)
			},
		)
	}
}

func TestShellLineHandlerJoinTwice(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				joined1, err := shelllinehandler.Join([]string{"echo", "hello \"world"})
				require.NoError(err)

				joined2, err := shelllinehandler.Join([]string{"bash", "-c", joined1})
				require.NoError(err)

				expected := "bash -c 'echo '\"'\"'hello \"world'\"'\"''"
				require.EqualValues(expected, joined2)

				for _, joined := range []string{joined1, joined2} {
					cmd := exec.Command("bash", "-c", joined)
					output, err := cmd.Output()
					require.NoError(err)
					executedOutput := strings.TrimSpace(string(output))
					require.EqualValues("hello \"world", executedOutput)
				}
			},
		)
	}
}

func TestShellLineHandlerJoinThreeTimes(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				joined1, err := shelllinehandler.Join([]string{"echo", "hello \"world"})
				require.NoError(err)

				joined2, err := shelllinehandler.Join([]string{"bash", "-c", joined1})
				require.NoError(err)

				joined3, err := shelllinehandler.Join([]string{"bash", "-c", joined2})
				require.NoError(err)

				expected := "bash -c 'bash -c '\"'\"'echo '\"'\"'\"'\"'\"'\"'\"'\"'hello \"world'\"'\"'\"'\"'\"'\"'\"'\"''\"'\"''"
				require.EqualValues(expected, joined3)

				for _, joined := range []string{joined1, joined2, joined3} {
					cmd := exec.Command("bash", "-c", joined)
					output, err := cmd.Output()
					require.NoError(err)
					executedOutput := strings.TrimSpace(string(output))
					require.EqualValues("hello \"world", executedOutput)
				}
			},
		)
	}
}
