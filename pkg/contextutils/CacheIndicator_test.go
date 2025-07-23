package contextutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
)


func Test_IsCacheIndicatorPresent(t *testing.T) {
	t.Run("not present", func(t *testing.T) {
		require.False(t, contextutils.IsCacheIndicatorPresent(context.Background()))
	})

	t.Run("present", func(t *testing.T) {
		require.True(t, contextutils.IsCacheIndicatorPresent(contextutils.WithCacheIndicator(context.Background())))
	})
}

func Test_CacheIndicator(t *testing.T) {
	t.Run("no cache indicator", func(t *testing.T) {
		ctx := context.Background()

		contextutils.SetCacheIndicator(ctx, true)

		isCached, err := contextutils.IsCachedResult(ctx)
		require.Error(t, err)
		require.False(t, isCached)
	})

	t.Run("cached result", func(t *testing.T) {
		ctx := contextutils.WithCacheIndicator(contextutils.ContextVerbose())

		isCached, err := contextutils.IsCachedResult(ctx)
		require.NoError(t, err)
		require.False(t, isCached)

		contextutils.SetCacheIndicator(ctx, true)

		isCached, err = contextutils.IsCachedResult(ctx)
		require.NoError(t, err)
		require.True(t, isCached)
	})
}
