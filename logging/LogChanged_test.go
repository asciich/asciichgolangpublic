package logging_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func Test_LogChangedSetsChangeIndicator(t *testing.T) {
	t.Run("testcase", func(t *testing.T) {
		ctx := contextutils.WithChangeIndicator(contextutils.ContextVerbose())

		// By default the verbose context does not indicate a change:
		require.False(t, contextutils.IsChanged(ctx))

		// Loging an info does not indicate a change
		logging.LogInfoByCtx(ctx, "only an info")
		logging.LogInfoByCtxf(ctx, "only an info: %d", 2)
		require.False(t, contextutils.IsChanged(ctx))

		// Loging an change does indicate a change
		logging.LogChangedByCtx(ctx, "log a change")
		require.True(t, contextutils.IsChanged(ctx))

		// Repeat again to vailidate LogChangedByCtxf 
		ctx = contextutils.WithChangeIndicator(contextutils.ContextVerbose())
		require.False(t, contextutils.IsChanged(ctx))
		logging.LogChangedByCtxf(ctx, "log a change: %d", 2)
		require.True(t, contextutils.IsChanged(ctx))
	})
}
