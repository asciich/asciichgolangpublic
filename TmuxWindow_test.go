package asciichgolangpublic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemuxWindow_CreateAndDeleteWindow(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Not available in Github CI")
		return
	}

	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				tmux := MustGetTmuxOnLocalMachine()

				session := tmux.MustGetSessionByName("sessionName")
				defer session.MustDelete(verbose)

				session.MustRecreate(verbose)

				window := session.MustGetWindowByName("windowName")

				for i := 0; i < 2; i++ {
					window.MustDelete(verbose)
					assert.False(window.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					window.MustCreate(verbose)
					assert.True(window.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					window.MustDelete(verbose)
					assert.False(window.MustExists(verbose))
				}
			},
		)
	}
}

func TestTemuxWindow_ReadLastLine(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogInfo("Not available in Github CI")
		return
	}

	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				tmux := MustGetTmuxOnLocalMachine()

				session := tmux.MustGetSessionByName("sessionName")
				defer session.MustDelete(verbose)

				session.MustRecreate(verbose)

				window := session.MustGetWindowByName("windowName")

				window.MustCreate(verbose)

				for i := 0; i < 3; i++ {
					content := RandomGenerator().MustGetRandomString(10)

					time.Sleep(2 * time.Second)

					window.MustSendKeys([]string{"echo " + content, "enter"}, verbose)

					time.Sleep(2 * time.Second)

					assert.EqualValues(
						content,
						window.MustGetSecondLatestPaneLine(),
					)
				}
			},
		)
	}
}
