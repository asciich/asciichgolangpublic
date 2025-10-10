package nativefiles_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestListFiles(t *testing.T) {
	t.Run("empty dir allow empty list", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		fileList, err := nativefiles.ListFiles(ctx, dir, &parameteroptions.ListFileOptions{
			AllowEmptyListIfNoFileIsFound: true,
		})
		require.NoError(t, err)
		require.Empty(t, fileList)
	})

	t.Run("empty dir", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		fileList, err := nativefiles.ListFiles(ctx, dir, &parameteroptions.ListFileOptions{
			AllowEmptyListIfNoFileIsFound: false,
		})
		require.Error(t, err)
		require.Nil(t, fileList)
	})
}

func Test_ListFilesWithSymlinks(t *testing.T) {
	t.Run("without symlink", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "a.txt")))
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "b.txt")))

		fileList, err := nativefiles.ListFiles(ctx, dir, &parameteroptions.ListFileOptions{})
		require.NoError(t, err)
		expected := []string{
			filepath.Join(dir, "a.txt"),
			filepath.Join(dir, "b.txt"),
		}
		require.EqualValues(t, expected, fileList)
	})

	t.Run("symlink to file", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		dir2, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir2, &filesoptions.DeleteOptions{}) }()

		// regular files
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "a.txt")))
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "b.txt")))

		// symlink
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir2, "target")))
		require.NoError(t, nativefiles.CreateSymlink(ctx, filepath.Join(dir2, "target"), filepath.Join(dir, "c.txt")))

		fileList, err := nativefiles.ListFiles(ctx, dir, &parameteroptions.ListFileOptions{})
		require.NoError(t, err)
		expected := []string{
			filepath.Join(dir, "a.txt"),
			filepath.Join(dir, "b.txt"),
			filepath.Join(dir, "c.txt"),
		}
		require.EqualValues(t, expected, fileList)
	})

	t.Run("symlink to directory", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		dir2, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir2, &filesoptions.DeleteOptions{}) }()

		// regular files
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "a.txt")))
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "b.txt")))
		require.NoError(t, nativefiles.CreateDirectory(ctx, filepath.Join(dir, "directory"), &filesoptions.CreateOptions{}))
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "directory", "c.txt")))
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir, "directory", "d.txt")))

		// symlink
		require.NoError(t, nativefiles.CreateDirectory(ctx, filepath.Join(dir2, "target"), &filesoptions.CreateOptions{}))
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir2, "target", "e.txt")))
		require.NoError(t, nativefiles.Create(ctx, filepath.Join(dir2, "target", "f.txt")))
		require.NoError(t, nativefiles.CreateSymlink(ctx, filepath.Join(dir2, "target"), filepath.Join(dir, "symlink")))

		fileList, err := nativefiles.ListFiles(ctx, dir, &parameteroptions.ListFileOptions{})
		require.NoError(t, err)
		expected := []string{
			filepath.Join(dir, "a.txt"),
			filepath.Join(dir, "b.txt"),
			filepath.Join(dir, "directory/c.txt"),
			filepath.Join(dir, "directory/d.txt"),
			filepath.Join(dir, "symlink/e.txt"),
			filepath.Join(dir, "symlink/f.txt"),
		}
		require.EqualValues(t, expected, fileList)
	})
}
