package nativefiles_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func Test_CreateAndDeleteSymlink(t *testing.T) {
	t.Run("create and delete for file", func(t *testing.T) {
		ctx := getCtx()

		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{}) }()

		targetFile := filepath.Join(tempDir, "target")
		require.NoError(t, nativefiles.Create(ctx, targetFile))

		symlink := filepath.Join(tempDir, "symlink")

		for range 2 {
			err = nativefiles.Delete(ctx, symlink, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.True(t, nativefiles.Exists(ctx, targetFile))
			require.False(t, nativefiles.Exists(ctx, symlink))
		}

		ctxChangeIndicator := contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetFile, symlink))
		require.True(t, nativefiles.Exists(ctx, targetFile))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.True(t, contextutils.IsChanged(ctxChangeIndicator))

		// Creating the symlink twice must be idempotent and not indicate a change:
		ctxChangeIndicator = contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetFile, symlink))
		require.True(t, nativefiles.Exists(ctx, targetFile))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.False(t, contextutils.IsChanged(ctxChangeIndicator))

		for range 2 {
			err = nativefiles.Delete(ctx, symlink, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.True(t, nativefiles.Exists(ctx, targetFile))
			require.False(t, nativefiles.Exists(ctx, symlink))
		}
	})

	t.Run("create and delete for dir", func(t *testing.T) {
		ctx := getCtx()

		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{}) }()

		targetDir := filepath.Join(tempDir, "target")
		require.NoError(t, nativefiles.CreateDirectory(ctx, targetDir))

		symlink := filepath.Join(tempDir, "symlink")

		for range 2 {
			err = nativefiles.Delete(ctx, symlink, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.True(t, nativefiles.Exists(ctx, targetDir))
			require.False(t, nativefiles.Exists(ctx, symlink))
		}

		ctxChangeIndicator := contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetDir, symlink))
		require.True(t, nativefiles.Exists(ctx, targetDir))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.True(t, contextutils.IsChanged(ctxChangeIndicator))

		// Creating the symlink twice must be idempotent and not indicate a change:
		ctxChangeIndicator = contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetDir, symlink))
		require.True(t, nativefiles.Exists(ctx, targetDir))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.False(t, contextutils.IsChanged(ctxChangeIndicator))

		for range 2 {
			err = nativefiles.Delete(ctx, symlink, &filesoptions.DeleteOptions{})
			require.NoError(t, err)

			require.True(t, nativefiles.Exists(ctx, targetDir))
			require.False(t, nativefiles.Exists(ctx, symlink))
		}
	})
}

func Test_CreateToAnotherTargetChangesTheSymlink(t *testing.T) {
	t.Run("create to changed target", func(t *testing.T) {
		ctx := getCtx()

		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{}) }()

		targetFile := filepath.Join(tempDir, "target")
		require.NoError(t, nativefiles.Create(ctx, targetFile))

		targetFile2 := filepath.Join(tempDir, "target2")
		require.NoError(t, nativefiles.Create(ctx, targetFile2))

		symlink := filepath.Join(tempDir, "symlink")

		require.True(t, nativefiles.Exists(ctx, targetFile))
		require.True(t, nativefiles.Exists(ctx, targetFile2))
		require.False(t, nativefiles.Exists(ctx, symlink))

		// Create to first target:
		ctxChangeIndicator := contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetFile, symlink))
		require.True(t, nativefiles.Exists(ctx, targetFile))
		require.True(t, nativefiles.Exists(ctx, targetFile2))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.True(t, contextutils.IsChanged(ctxChangeIndicator))
		require.True(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile)))
		require.False(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile2)))

		// Create to first target again to check idempotence:
		ctxChangeIndicator = contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetFile, symlink))
		require.True(t, nativefiles.Exists(ctx, targetFile))
		require.True(t, nativefiles.Exists(ctx, targetFile2))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.False(t, contextutils.IsChanged(ctxChangeIndicator))
		require.True(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile)))
		require.False(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile2)))

		// Create to second target:
		ctxChangeIndicator = contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetFile2, symlink))
		require.True(t, nativefiles.Exists(ctx, targetFile))
		require.True(t, nativefiles.Exists(ctx, targetFile2))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.True(t, contextutils.IsChanged(ctxChangeIndicator))
		require.False(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile)))
		require.True(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile2)))

		// Create to second target again to check idempotence:
		ctxChangeIndicator = contextutils.WithChangeIndicator(ctx)
		require.NoError(t, nativefiles.CreateSymlink(ctxChangeIndicator, targetFile2, symlink))
		require.True(t, nativefiles.Exists(ctx, targetFile))
		require.True(t, nativefiles.Exists(ctx, targetFile2))
		require.True(t, nativefiles.Exists(ctx, symlink))
		require.False(t, contextutils.IsChanged(ctxChangeIndicator))
		require.False(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile)))
		require.True(t, mustutils.Must(nativefiles.IsSymlinkTo(ctx, symlink, targetFile2)))
	})
}

func Test_IsSymlink(t *testing.T) {
	t.Run("empty string as path", func(t *testing.T) {
		ctx := getCtx()

		isSymlink, err := nativefiles.IsSymlink(ctx, "")
		require.Error(t, err)
		require.False(t, isSymlink)
	})

	t.Run("non existing file", func(t *testing.T) {
		ctx := getCtx()

		isSymlink, err := nativefiles.IsSymlink(ctx, "/this_file_does_not_exist")
		require.Error(t, err)
		require.False(t, isSymlink)
	})

	t.Run("dir is not a symlink", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		isSymlink, err := nativefiles.IsSymlink(ctx, dir)
		require.NoError(t, err)
		require.False(t, isSymlink)
	})

	t.Run("regular file is not a symlink", func(t *testing.T) {
		ctx := getCtx()

		file, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, file, &filesoptions.DeleteOptions{}) }()

		isSymlink, err := nativefiles.IsSymlink(ctx, file)
		require.NoError(t, err)
		require.False(t, isSymlink)
	})

	t.Run("symlink to file is symlink", func(t *testing.T) {
		ctx := getCtx()

		file, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, file, &filesoptions.DeleteOptions{}) }()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		symlinkPath := filepath.Join(dir, "symlink")

		err = nativefiles.CreateSymlink(ctx, file, symlinkPath)
		require.NoError(t, err)

		isSymlink, err := nativefiles.IsSymlink(ctx, symlinkPath)
		require.NoError(t, err)
		require.True(t, isSymlink)
	})
}


func Test_IsSymlinkToDirectory(t *testing.T) {
	t.Run("empty string as path", func(t *testing.T) {
		ctx := getCtx()

		isSymlink, err := nativefiles.IsSymlinkToDirectory(ctx, "")
		require.Error(t, err)
		require.False(t, isSymlink)
	})

	t.Run("non existing file", func(t *testing.T) {
		ctx := getCtx()

		isSymlink, err := nativefiles.IsSymlinkToDirectory(ctx, "/this_file_does_not_exist")
		require.Error(t, err)
		require.False(t, isSymlink)
	})

	t.Run("symlink to file", func(t *testing.T) {
		ctx := getCtx()

		file, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, file, &filesoptions.DeleteOptions{}) }()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		symlinkPath := filepath.Join(dir, "symlink")

		err = nativefiles.CreateSymlink(ctx, file, symlinkPath)
		require.NoError(t, err)

		isSymlinkToDir, err := nativefiles.IsSymlinkToDirectory(ctx, symlinkPath)
		require.NoError(t, err)
		require.False(t, isSymlinkToDir)
	})


	t.Run("symlink to directory", func(t *testing.T) {
		ctx := getCtx()

		targetDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, targetDir, &filesoptions.DeleteOptions{}) }()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, dir, &filesoptions.DeleteOptions{}) }()

		symlinkPath := filepath.Join(dir, "symlink")

		err = nativefiles.CreateSymlink(ctx, targetDir, symlinkPath)
		require.NoError(t, err)

		isSymlinkToDir, err := nativefiles.IsSymlinkToDirectory(ctx, symlinkPath)
		require.NoError(t, err)
		require.True(t, isSymlinkToDir)
	})
}
