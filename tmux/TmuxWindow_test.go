package tmux

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/randomgenerator"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestTemuxWindow_CreateAndDeleteWindow(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Not available in Github CI")
		return
	}

	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				tmux := MustGetTmuxOnLocalMachine()

				session := tmux.MustGetSessionByName("sessionName")
				defer session.MustDelete(verbose)

				session.MustRecreate(verbose)

				window := session.MustGetWindowByName("windowName")

				for i := 0; i < 2; i++ {
					window.MustDelete(verbose)
					require.False(window.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					window.MustCreate(verbose)
					require.True(window.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					window.MustDelete(verbose)
					require.False(window.MustExists(verbose))
				}
			},
		)
	}
}

func TestTemuxWindow_ReadLastLine(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Not available in Github CI")
		return
	}

	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				tmux := MustGetTmuxOnLocalMachine()

				session := tmux.MustGetSessionByName("sessionName")
				defer session.MustDelete(verbose)

				session.MustRecreate(verbose)

				window := session.MustGetWindowByName("windowName")

				window.MustCreate(verbose)

				for i := 0; i < 3; i++ {
					content := randomgenerator.MustGetRandomString(10)

					window.MustWaitUntilCliPromptReady(verbose)

					window.MustSendKeys([]string{"echo " + content, "enter"}, verbose)

					time.Sleep(time.Millisecond * 500)

					require.EqualValues(
						content,
						window.MustGetSecondLatestPaneLine(),
					)
				}
			},
		)
	}
}

func TestTemuxWindow_RunCommand(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Not available in Github CI")
		return
	}

	tests := []struct {
		command        []string
		expectedStdout string
	}{
		{[]string{"echo", "hello"}, "hello\n"},
		{[]string{"bash", "-c", "echo hello"}, "hello\n"},
		{[]string{"bash", "-c", "sleep 2s ; echo hello"}, "hello\n"},
		{[]string{"echo", "-en", "hello"}, "hello"},
		{[]string{"echo", "hello world"}, "hello world\n"},
		{[]string{"echo", "hello", "world"}, "hello world\n"},
		{[]string{"echo", "-en", "hello\\nworld\\n"}, "hello\nworld\n"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				tmux := MustGetTmuxOnLocalMachine()

				window := tmux.MustGetWindowByNames("sessionName", "windowName")
				defer window.MustDeleteSession(verbose)

				window.MustRecreate(verbose)

				commandOutput := window.MustRunCommand(
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)

				require.EqualValues(
					tt.expectedStdout,
					commandOutput.MustGetStdoutAsString(),
				)
			},
		)
	}
}
