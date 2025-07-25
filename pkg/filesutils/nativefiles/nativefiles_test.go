package nativefiles_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_IsFile(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		require.False(t, nativefiles.IsFile(ctx, ""))
	})

	t.Run("nonexisting file", func(t *testing.T) {
		require.False(t, nativefiles.IsFile(ctx, "/this/file/does/not/exist"))
	})

	t.Run("existing file", func(t *testing.T) {
		require.True(t, nativefiles.IsFile(ctx, "/etc/hosts"))
	})

	t.Run("directory", func(t *testing.T) {
		require.False(t, nativefiles.IsFile(ctx, "/etc"))
	})

	t.Run("directory2", func(t *testing.T) {
		require.False(t, nativefiles.IsFile(ctx, "/etc/"))
	})
}

func Test_IsDir(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		require.False(t, nativefiles.IsDir(ctx, ""))
	})

	t.Run("nonexisting file", func(t *testing.T) {
		require.False(t, nativefiles.IsDir(ctx, "/this/file/does/not/exist"))
	})

	t.Run("existing file", func(t *testing.T) {
		require.False(t, nativefiles.IsDir(ctx, "/etc/hosts"))
	})

	t.Run("directory", func(t *testing.T) {
		require.True(t, nativefiles.IsDir(ctx, "/etc"))
	})

	t.Run("directory2", func(t *testing.T) {
		require.True(t, nativefiles.IsDir(ctx, "/etc/"))
	})
}

func Test_Exists(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		require.False(t, nativefiles.Exists(ctx, ""))
	})

	t.Run("nonexisting file", func(t *testing.T) {
		require.False(t, nativefiles.Exists(ctx, "/this/file/does/not/exist"))
	})

	t.Run("existing file", func(t *testing.T) {
		require.True(t, nativefiles.Exists(ctx, "/etc/hosts"))
	})

	t.Run("directory", func(t *testing.T) {
		require.True(t, nativefiles.Exists(ctx, "/etc"))
	})

	t.Run("directory2", func(t *testing.T) {
		require.True(t, nativefiles.Exists(ctx, "/etc/"))
	})
}

func Test_Create(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		err := nativefiles.Create(ctx, "")
		require.Error(t, err)
	})

	t.Run("create", func(t *testing.T) {
		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)

		testPath := filepath.Join(tempDir, "test.txt")
		require.False(t, nativefiles.IsFile(ctx, testPath))

		ctx := contextutils.WithChangeIndicator(ctx)
		err = nativefiles.Create(ctx, testPath)
		require.NoError(t, err)
		require.True(t, nativefiles.IsFile(ctx, testPath))
		require.True(t, contextutils.IsChanged(ctx))

		ctx = contextutils.WithChangeIndicator(ctx)
		err = nativefiles.Create(ctx, testPath)
		require.NoError(t, err)
		require.True(t, nativefiles.IsFile(ctx, testPath))
		require.False(t, contextutils.IsChanged(ctx))
	})
}
