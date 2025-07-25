package tmuxutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/tmuxutils"
)

func TestTemuxWindow_CreateAndDeleteWindow(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tmux, err := tmuxutils.GetTmuxOnLocalMachine()
				require.NoError(t, err)

				session, err := tmux.GetSessionByName("sessionName")
				require.NoError(t, err)
				defer session.Delete(ctx)

				err = session.Recreate(ctx)
				require.NoError(t, err)

				window, err := session.GetWindowByName("windowName")
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = window.Delete(ctx)
					require.NoError(t, err)

					windowExists, err := window.Exists(ctx)
					require.NoError(t, err)
					require.False(t, windowExists)
				}

				for i := 0; i < 2; i++ {
					err = window.Create(ctx)
					require.NoError(t, err)

					windowExists, err := window.Exists(ctx)
					require.NoError(t, err)
					require.True(t, windowExists)
				}

				for i := 0; i < 2; i++ {
					err = window.Delete(ctx)
					require.NoError(t, err)

					windowExists, err := window.Exists(ctx)
					require.NoError(t, err)
					require.False(t, windowExists)
				}
			},
		)
	}
}

func TestTemuxWindow_ReadLastLine(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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
				ctx := getCtx()

				tmux, err := tmuxutils.GetTmuxOnLocalMachine()
				require.NoError(t, err)

				session, err := tmux.GetSessionByName("sessionName")
				require.NoError(t, err)
				defer session.Delete(ctx)

				err = session.Recreate(ctx)
				require.NoError(t, err)

				window, err := session.GetWindowByName("windowName")
				require.NoError(t, err)

				err = window.Create(ctx)
				require.NoError(t, err)

				err = window.WaitUntilCliPromptReady(ctx)
				require.NoError(t, err)

				err = window.SendKeys(ctx, []string{"echo '" + tt.testmessage + "'", "enter"})
				require.NoError(t, err)

				err = window.WaitUntilCliPromptReady(ctx)
				require.NoError(t, err)

				line, err := window.GetSecondLatestPaneLine()
				require.NoError(t, err)
				require.EqualValues(
					t,
					tt.testmessage,
					line,
				)
			},
		)
	}
}

func TestTemuxWindow_WaitOutputMatchesRegex(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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
				ctx := getCtx()

				tmux, err := tmuxutils.GetTmuxOnLocalMachine()
				require.NoError(t, err)

				session, err := tmux.GetSessionByName("sessionName")
				require.NoError(t, err)
				defer session.Delete(ctx)

				err = session.Recreate(ctx)
				require.NoError(t, err)

				window, err := session.GetWindowByName("windowName")
				require.NoError(t, err)

				err = window.Create(ctx)
				require.NoError(t, err)

				err = window.WaitUntilCliPromptReady(ctx)
				require.NoError(t, err)

				outputPath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(contextutils.GetVerboseFromContext(ctx))
				require.NoError(t, err)
				defer files.DeleteFileByPath(outputPath, contextutils.GetVerboseFromContext(ctx))

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

				exampleScriptPath, err := tempfilesoo.CreateFromStringAndGetPath(exampleScript, contextutils.GetVerboseFromContext(ctx))
				require.NoError(t, err)
				defer files.DeleteFileByPath(exampleScriptPath, contextutils.GetVerboseFromContext(ctx))

				err = window.SendKeys(
					ctx,
					[]string{
						"bash " + exampleScriptPath,
						"enter",
					},
				)
				require.NoError(t, err)

				err = window.WaitUntilOutputMatchesRegex(ctx, "Username:", 2*time.Second)
				require.NoError(t, err)
				err = window.SendKeys(ctx, []string{tt.username, "enter"})
				require.NoError(t, err)
				err = window.WaitUntilOutputMatchesRegex(ctx, "Password:", 2*time.Second)
				require.NoError(t, err)
				err = window.SendKeys(ctx, []string{tt.password, "enter"})
				require.NoError(t, err)
				err = window.WaitUntilOutputMatchesRegex(ctx, "finished", 2*time.Second)
				require.NoError(t, err)

				shownLines, err := window.GetShownLines()
				require.NoError(t, err)

				content, err := files.ReadFileAsString(outputPath)
				require.NoError(t, err)
				require.EqualValues(t, tt.username+"\n"+tt.password+"\n", content)
				require.Contains(t, shownLines, "finished")
			},
		)
	}
}

func TestTemuxWindow_RunCommand(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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
				ctx := getCtx()

				tmux, err := tmuxutils.GetTmuxOnLocalMachine()
				require.NoError(t, err)

				window, err := tmux.GetWindowByNames("sessionName", "windowName")
				require.NoError(t, err)
				defer window.DeleteSession(ctx)

				err = window.Recreate(ctx)
				require.NoError(t, err)

				commandOutput, err := window.RunCommand(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: tt.command,
					},
				)
				require.NoError(t, err)

				stdout, err := commandOutput.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedStdout, stdout)
			},
		)
	}
}
