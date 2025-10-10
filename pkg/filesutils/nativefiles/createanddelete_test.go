package nativefiles_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_CreateAndDeleteFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		// We use a temporary file path for testing:
		filePath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
		err = nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
		require.NoError(t, err)

		// Test Delete first to ensure the file is absent in any case for the next test steps.
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.False(t, nativefiles.Exists(ctx, filePath))
		}

		// Create the file again
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Create(ctx, filePath)
			require.NoError(t, err)

			require.True(t, nativefiles.Exists(ctx, filePath))
		}

		// Test Delete.
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.False(t, nativefiles.Exists(ctx, filePath))
		}
	})
}

func Test_CreateAndDeleteDir(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		// We use a temporary directory path for testing:
		filePath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
		err = nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
		require.NoError(t, err)

		// Test Delete first to ensure the file is absent in any case for the next test steps.
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.False(t, nativefiles.Exists(ctx, filePath))
		}

		// Create the file again
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.CreateDirectory(ctx, filePath, &filesoptions.CreateOptions{})
			require.NoError(t, err)

			require.True(t, nativefiles.Exists(ctx, filePath))
			require.True(t, nativefiles.IsDir(ctx, filePath))
		}

		// Test Delete.
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.False(t, nativefiles.Exists(ctx, filePath))
		}
	})
}

func Test_DeleteDirRecursively(t *testing.T) {
	t.Run("delete empty dir", func(t *testing.T) {
		ctx := getCtx()
		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)

		require.DirExists(t, tempDir)
		err = nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})
		require.NoError(t, err)

		require.NoDirExists(t, tempDir)
	})

	t.Run("delete non empty dir", func(t *testing.T) {
		ctx := getCtx()
		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)

		require.DirExists(t, tempDir)

		err = nativefiles.CreateDirectory(ctx, filepath.Join(tempDir, "abc"), &filesoptions.CreateOptions{})
		require.NoError(t, err)

		err = nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})
		require.NoError(t, err)

		require.NoDirExists(t, tempDir)
	})
}
