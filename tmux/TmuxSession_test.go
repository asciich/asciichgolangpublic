package tmux

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestTemuxSession_CreateAndDeleteSession(t *testing.T) {
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
				require := require.New(t)

				const verbose bool = true

				tmux := MustGetTmuxOnLocalMachine()

				session := tmux.MustGetSessionByName("sessionName")
				defer session.MustDelete(verbose)

				for i := 0; i < 2; i++ {
					session.MustDelete(verbose)
					require.False(session.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					session.MustCreate(verbose)
					require.True(session.MustExists(verbose))
				}

				time.Sleep(1 * time.Second)

				for i := 0; i < 2; i++ {
					session.MustDelete(verbose)
					require.False(session.MustExists(verbose))
				}
			},
		)
	}
}
