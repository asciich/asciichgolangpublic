package nativefiles_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func TestSecureDelete(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		err := nativefiles.SecureDelete(getCtx(), "")
		require.Error(t, err)
	})

	t.Run("non existing path", func(t *testing.T) {
		ctx := contextutils.WithChangeIndicator(getCtx())
		err := nativefiles.SecureDelete(ctx, "/this/path/does/not/exist")
		require.NoError(t, err)
		require.False(t, contextutils.IsChanged(ctx))
	})

	t.Run("non existing path", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		require.True(t, nativefiles.IsFile(ctx, tempPath))

		ctx2 := contextutils.WithChangeIndicator(ctx)
		err = nativefiles.SecureDelete(ctx2, tempPath)
		require.NoError(t, err)
		require.True(t, contextutils.IsChanged(ctx2))
		require.False(t, nativefiles.IsFile(ctx, tempPath))

		ctx2 = contextutils.WithChangeIndicator(ctx)
		err = nativefiles.SecureDelete(ctx2, tempPath)
		require.NoError(t, err)
		require.False(t, contextutils.IsChanged(ctx2))
		require.False(t, nativefiles.IsFile(ctx, tempPath))
	})
}
