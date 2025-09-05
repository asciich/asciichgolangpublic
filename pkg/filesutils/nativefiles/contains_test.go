package nativefiles_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func TestContains(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})}()

		contains, err := nativefiles.Contains(ctx, tempFile, "hello world")
		require.NoError(t, err)
		require.False(t, contains)
	})

	t.Run("contains", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})}()

		contains, err := nativefiles.Contains(ctx, tempFile, "hello")
		require.NoError(t, err)
		require.True(t, contains)
	})

	t.Run("not contains", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer func() { _ = nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})}()

		contains, err := nativefiles.Contains(ctx, tempFile, "not included")
		require.NoError(t, err)
		require.False(t, contains)
	})
}
