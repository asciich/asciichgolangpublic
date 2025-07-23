package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/filesutils"
)

func TestSecureDelete(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		err := filesutils.SecureDelete(getCtx(), "")
		require.Error(t, err)
	})

	t.Run("non existing path", func(t *testing.T) {
		ctx := contextutils.WithChangeIndicator(getCtx())
		err := filesutils.SecureDelete(ctx, "/this/path/does/not/exist")
		require.NoError(t, err)
		require.False(t, contextutils.IsChanged(ctx))
	})

	t.Run("non existing path", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := filesutils.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		require.True(t, filesutils.IsFile(ctx, tempPath))

		ctx2 := contextutils.WithChangeIndicator(ctx)
		err = filesutils.SecureDelete(ctx2, tempPath)
		require.NoError(t, err)
		require.True(t, contextutils.IsChanged(ctx2))
		require.False(t, filesutils.IsFile(ctx, tempPath))

		ctx2 = contextutils.WithChangeIndicator(ctx)
		err = filesutils.SecureDelete(ctx2, tempPath)
		require.NoError(t, err)
		require.False(t, contextutils.IsChanged(ctx2))
		require.False(t, filesutils.IsFile(ctx, tempPath))
	})
}
