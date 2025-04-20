package tmux

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

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
				ctx := getCtx()

				tmux := MustGetTmuxOnLocalMachine()

				session := tmux.MustGetSessionByName("sessionName")
				defer session.Delete(ctx)

				for i := 0; i < 2; i++ {
					err := session.Delete(ctx)
					require.NoError(t, err)

					exists, err := session.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}

				for i := 0; i < 2; i++ {
					err := session.Create(ctx)
					require.NoError(t, err)

					exists, err := session.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)
				}

				time.Sleep(1 * time.Second)

				for i := 0; i < 2; i++ {
					err := session.Delete(ctx)
					require.NoError(t, err)

					exists, err := session.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
