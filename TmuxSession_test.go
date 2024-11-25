package asciichgolangpublic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemuxSession_CreateAndDeleteSession(t *testing.T) {
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

				for i := 0; i < 2; i++ {
					session.MustDelete(verbose)
					assert.False(session.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					session.MustCreate(verbose)
					assert.True(session.MustExists(verbose))
				}

				time.Sleep(1 * time.Second)

				for i := 0; i < 2; i++ {
					session.MustDelete(verbose)
					assert.False(session.MustExists(verbose))
				}
			},
		)
	}
}
