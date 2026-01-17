package commandexecutorfileoo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutortempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateDirectory(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		dir, err := commandexecutortempfilesoo.CreateLocalEmptyTemporaryDirectory(ctx)
		require.NoError(t, err)
		defer dir.Delete(ctx, &filesoptions.DeleteOptions{})

		require.NoError(t, dir.Delete(ctx, &filesoptions.DeleteOptions{}))

		exists, err := dir.Exists(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		require.NoError(t, dir.Create(ctx, &filesoptions.CreateOptions{}))

		exists, err = dir.Exists(ctx)
		require.NoError(t, err)
		require.True(t, exists)

		path, err := dir.GetPath()
		require.NoError(t, err)
		require.True(t, nativefiles.IsDir(ctx, path))
	})
}

func Test_CreateFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		file, err := commandexecutortempfilesoo.CreateLocalEmptyTemporaryFile(ctx)
		require.NoError(t, err)
		defer file.Delete(ctx, &filesoptions.DeleteOptions{})

		require.NoError(t, file.Delete(ctx, &filesoptions.DeleteOptions{}))

		exists, err := file.Exists(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		require.NoError(t, file.Create(ctx, &filesoptions.CreateOptions{}))

		exists, err = file.Exists(ctx)
		require.NoError(t, err)
		require.True(t, exists)

		path, err := file.GetPath()
		require.NoError(t, err)
		require.True(t, nativefiles.IsFile(ctx, path))
	})
}
