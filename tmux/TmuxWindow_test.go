package tmux

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tempfiles"
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
		{"Aengia0s"},
		{"Gu2aivai"},
		{"Aen8ayai"},
		{"Aen8a;yai"},
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

				window.MustWaitUntilCliPromptReady(verbose)

				window.MustSendKeys([]string{"echo '" + tt.testmessage + "'", "enter"}, verbose)

				window.MustWaitUntilCliPromptReady(verbose)

				require.EqualValues(
					tt.testmessage,
					window.MustGetSecondLatestPaneLine(),
				)
			},
		)
	}
}

func TestTemuxWindow_WaitOutputMatchesRegex(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Not available in Github CI")
		return
	}

	tests := []struct {
		username string
		password string
	}{
		{"user1", "Aengia0s"},
		{"user2", "Aengsdfsdfa0s"},
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

				window.MustWaitUntilCliPromptReady(verbose)

				outputPath := tempfiles.MustCreateEmptyTemporaryFileAndGetPath(verbose)
				defer files.MustDeleteFileByPath(outputPath, verbose)

				exampleScript := "#/usr/bin/env bash\n"
				exampleScript += "\n"
				exampleScript += "sleep 0.5\n"
				exampleScript += "echo Username:\n"
				exampleScript += "read USERNAME\n"
				exampleScript += "echo $USERNAME >> '" + outputPath + "'\n"
				exampleScript += "sleep 1\n"
				exampleScript += "echo Password:\n"
				exampleScript += "read PASSWORD\n"
				exampleScript += "echo $PASSWORD >> '" + outputPath + "'\n"
				exampleScript += "sleep .75\n"
				exampleScript += "echo finished\n"

				exampleScriptPath := tempfiles.MustCreateFromStringAndGetPath(exampleScript, verbose)
				defer files.MustDeleteFileByPath(exampleScriptPath, verbose)

				window.MustSendKeys(
					[]string{
						"bash " + exampleScriptPath,
						"enter",
					},
					verbose,
				)

				window.MustWaitUntilOutputMatchesRegex("Username:", 2*time.Second, verbose)
				window.MustSendKeys([]string{tt.username, "enter"}, verbose)
				window.MustWaitUntilOutputMatchesRegex("Password:", 2*time.Second, verbose)
				window.MustSendKeys([]string{tt.password, "enter"}, verbose)
				window.MustWaitUntilOutputMatchesRegex("finished", 2*time.Second, verbose)

				shownLines := window.MustGetShownLines()
				require.EqualValues(tt.username+"\n"+tt.password+"\n", files.MustReadFileAsString(outputPath))
				require.Contains(shownLines, "finished")
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
