package nativefiles_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func TestSetBlockInFile_BlockAdded(t *testing.T) {
	ctx := getCtx()

	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "some existing content\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	ctx = contextutils.WithChangeIndicator(ctx)
	err = nativefiles.SetBlockInFile(ctx, path, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)

	result, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.Contains(t, result, "# BEGIN EXAMPLE_BLOCK")
	require.Contains(t, result, "example content")
	require.Contains(t, result, "# END EXAMPLE_BLOCK")
	require.Contains(t, result, "some existing content")
	require.True(t, contextutils.IsChanged(ctx))
}

func TestSetBlockInFile_BlockAlreadyPresent(t *testing.T) {
	ctx := getCtx()

	content := "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK\n"
	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	ctx = contextutils.WithChangeIndicator(ctx)
	err = nativefiles.SetBlockInFile(ctx, path, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)

	result, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.EqualValues(t, content, result)
	require.False(t, contextutils.IsChanged(ctx))
}

func TestSetBlockInFile_BlockModified(t *testing.T) {
	ctx := getCtx()

	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "# BEGIN EXAMPLE_BLOCK\nold content\n# END EXAMPLE_BLOCK\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	ctx = contextutils.WithChangeIndicator(ctx)
	err = nativefiles.SetBlockInFile(ctx, path, "EXAMPLE_BLOCK", "new content")
	require.NoError(t, err)

	result, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.Contains(t, result, "# BEGIN EXAMPLE_BLOCK")
	require.Contains(t, result, "new content")
	require.Contains(t, result, "# END EXAMPLE_BLOCK")
	require.NotContains(t, result, "old content")
	require.True(t, contextutils.IsChanged(ctx))
}

func TestSetBlockInFile_ModifiedChangeIndicatorIsSet(t *testing.T) {
	ctx := getCtx()

	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "# BEGIN EXAMPLE_BLOCK\nold content\n# END EXAMPLE_BLOCK\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	ctx = contextutils.WithChangeIndicator(ctx)
	err = nativefiles.SetBlockInFile(ctx, path, "EXAMPLE_BLOCK", "new content")
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctx))
}

func TestSetBlockInFile_UnchangedChangeIndicatorNotSet(t *testing.T) {
	ctx := getCtx()

	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	ctx = contextutils.WithChangeIndicator(ctx)
	err = nativefiles.SetBlockInFile(ctx, path, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctx))
}

func TestSetBlockInFile_AddedToContentWithTrailingNewline(t *testing.T) {
	ctx := getCtx()

	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "some existing content\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	ctx = contextutils.WithChangeIndicator(ctx)
	err = nativefiles.SetBlockInFile(ctx, path, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)

	result, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctx))
	require.Contains(t, result, "some existing content\n")
	require.Contains(t, result, "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK")
}

func TestSetBlockInFile_AddedToContentWithoutTrailingNewline(t *testing.T) {
	ctx := getCtx()

	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "some existing content")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	ctx = contextutils.WithChangeIndicator(ctx)
	err = nativefiles.SetBlockInFile(ctx, path, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)

	result, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctx))
	require.Contains(t, result, "some existing content\n")
	require.Contains(t, result, "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK")
}

func TestSetBlockInFile_EmptyPathReturnsError(t *testing.T) {
	ctx := getCtx()

	err := nativefiles.SetBlockInFile(ctx, "", "EXAMPLE_BLOCK", "example content")
	require.Error(t, err)
}

func TestSetBlockInFile_EmptyBlockNameReturnsError(t *testing.T) {
	ctx := getCtx()

	path, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "some content\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, path, &filesoptions.DeleteOptions{})

	err = nativefiles.SetBlockInFile(ctx, path, "", "example content")
	require.Error(t, err)
}
