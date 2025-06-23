package contextutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func Test_IsChangeIndicatorPresent(t *testing.T) {
	t.Run("not present", func(t *testing.T) {
		require.False(t, contextutils.IsChangeIndicatorPresent(context.Background()))
	})

	t.Run("present", func(t *testing.T) {
		require.True(t, contextutils.IsChangeIndicatorPresent(contextutils.WithChangeIndicator(context.Background())))
	})
}

func Test_ChangeIndicator(t *testing.T) {
	t.Run("no Change indicator", func(t *testing.T) {
		ctx := context.Background()

		contextutils.SetChangeIndicator(ctx, true)

		isChanged, err := contextutils.IsChangedResult(ctx)
		require.Error(t, err)
		require.False(t, isChanged)

		require.False(t, contextutils.IsChanged(ctx))
	})

	t.Run("Changed result", func(t *testing.T) {
		ctx := contextutils.WithChangeIndicator(contextutils.ContextVerbose())

		isChanged, err := contextutils.IsChangedResult(ctx)
		require.NoError(t, err)
		require.False(t, isChanged)

		contextutils.SetChangeIndicator(ctx, true)

		isChanged, err = contextutils.IsChangedResult(ctx)
		require.NoError(t, err)
		require.True(t, isChanged)

		require.True(t, contextutils.IsChanged(ctx))
	})
}
