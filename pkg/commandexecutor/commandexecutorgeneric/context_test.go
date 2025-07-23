package commandexecutorgeneric_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func Test_WithLiveOutputOnStdout(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(nil))
	})

	t.Run("with nil", func(t *testing.T) {
		ctx := commandexecutorgeneric.WithLiveOutputOnStdout(nil)
		require.True(t, commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(ctx))
	})

	t.Run("with silent", func(t *testing.T) {
		ctx := commandexecutorgeneric.WithLiveOutputOnStdout(contextutils.ContextSilent())
		require.True(t, commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(ctx))
	})

	t.Run("with silent", func(t *testing.T) {
		ctx := commandexecutorgeneric.WithLiveOutputOnStdout(contextutils.ContextVerbose())
		require.True(t, commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(ctx))
	})

	t.Run("silent context", func(t *testing.T) {
		require.False(t, commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(contextutils.ContextSilent()))
	})

	t.Run("verbose context", func(t *testing.T) {
		require.False(t, commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(contextutils.ContextVerbose()))
	})

	t.Run("if verbose", func(t *testing.T) {
		require.False(
			t,
			commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(
				commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(contextutils.ContextSilent()),
			),
		)

		require.True(
			t,
			commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(
				commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(contextutils.ContextVerbose()),
			),
		)
	})
}
