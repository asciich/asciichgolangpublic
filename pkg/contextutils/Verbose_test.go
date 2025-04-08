package contextutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func Test_GetVerbosityContextByBool(t *testing.T) {
	t.Run("verbose", func(t *testing.T) {
		ctx := contextutils.GetVerbosityContextByBool(true)
		require.True(t, contextutils.GetVerboseFromContext(ctx))
	})

	t.Run("silent", func(t *testing.T) {
		ctx := contextutils.GetVerbosityContextByBool(false)
		require.False(t, contextutils.GetVerboseFromContext(ctx))
	})
}

func Test_DefaultContexts(t *testing.T) {
	t.Run("ContextVerbose", func(t *testing.T) {
		ctx := contextutils.ContextVerbose()
		require.True(t, contextutils.GetVerboseFromContext(ctx))
	})

	t.Run("ContextSilent", func(t *testing.T) {
		ctx := contextutils.ContextSilent()
		require.False(t, contextutils.GetVerboseFromContext(ctx))
	})
}
