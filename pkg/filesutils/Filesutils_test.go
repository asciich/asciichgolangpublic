package filesutils_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/filesutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_IsFile(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, ""))
	})

	t.Run("nonexisting file", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, "/this/file/does/not/exist"))
	})

	t.Run("existing file", func(t *testing.T) {
		require.True(t, filesutils.IsFile(ctx, "/etc/hosts"))
	})

	t.Run("directory", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, "/etc"))
	})

	t.Run("directory2", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, "/etc/"))
	})
}

func Test_IsDir(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		require.False(t, filesutils.IsDir(ctx, ""))
	})

	t.Run("nonexisting file", func(t *testing.T) {
		require.False(t, filesutils.IsDir(ctx, "/this/file/does/not/exist"))
	})

	t.Run("existing file", func(t *testing.T) {
		require.False(t, filesutils.IsDir(ctx, "/etc/hosts"))
	})

	t.Run("directory", func(t *testing.T) {
		require.True(t, filesutils.IsDir(ctx, "/etc"))
	})

	t.Run("directory2", func(t *testing.T) {
		require.True(t, filesutils.IsDir(ctx, "/etc/"))
	})
}

func Test_Exists(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		require.False(t, filesutils.Exists(ctx, ""))
	})

	t.Run("nonexisting file", func(t *testing.T) {
		require.False(t, filesutils.Exists(ctx, "/this/file/does/not/exist"))
	})

	t.Run("existing file", func(t *testing.T) {
		require.True(t, filesutils.Exists(ctx, "/etc/hosts"))
	})

	t.Run("directory", func(t *testing.T) {
		require.True(t, filesutils.Exists(ctx, "/etc"))
	})

	t.Run("directory2", func(t *testing.T) {
		require.True(t, filesutils.Exists(ctx, "/etc/"))
	})
}

func Test_Create(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		err := filesutils.Create(ctx, "")
		require.Error(t, err)
	})

	t.Run("create", func(t *testing.T) {
		tempDir, err := filesutils.CreateTempDir(ctx)
		require.NoError(t, err)

		testPath := filepath.Join(tempDir, "test.txt")
		require.False(t, filesutils.IsFile(ctx, testPath))

		ctx := contextutils.WithChangeIndicator(ctx)
		err = filesutils.Create(ctx, testPath)
		require.NoError(t, err)
		require.True(t, filesutils.IsFile(ctx, testPath))
		require.True(t, contextutils.IsChanged(ctx))

		ctx = contextutils.WithChangeIndicator(ctx)
		err = filesutils.Create(ctx, testPath)
		require.NoError(t, err)
		require.True(t, filesutils.IsFile(ctx, testPath))
		require.False(t, contextutils.IsChanged(ctx))
	})
}
