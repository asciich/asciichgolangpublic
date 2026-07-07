package stringsutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestBlockAlreadyPresent(t *testing.T) {
	ctx := getCtx()
	ctx = contextutils.WithChangeIndicator(ctx)

	content := "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK\n"

	updatedContent, err := stringsutils.BlockInString(ctx, content, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)
	require.EqualValues(t, content, updatedContent)
	require.False(t, contextutils.IsChanged(ctx))
}

func TestBlockAdded(t *testing.T) {
	ctx := getCtx()
	ctx = contextutils.WithChangeIndicator(ctx)

	content := "some existing content\n"

	updatedContent, err := stringsutils.BlockInString(ctx, content, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)
	require.Contains(t, updatedContent, "# BEGIN EXAMPLE_BLOCK")
	require.Contains(t, updatedContent, "example content")
	require.Contains(t, updatedContent, "# END EXAMPLE_BLOCK")
	require.Contains(t, updatedContent, "some existing content")
	require.True(t, contextutils.IsChanged(ctx))
}

func TestBlockModified(t *testing.T) {
	ctx := getCtx()
	ctx = contextutils.WithChangeIndicator(ctx)

	content := "# BEGIN EXAMPLE_BLOCK\nold content\n# END EXAMPLE_BLOCK\n"

	updatedContent, err := stringsutils.BlockInString(ctx, content, "EXAMPLE_BLOCK", "new content")
	require.NoError(t, err)
	require.Contains(t, updatedContent, "# BEGIN EXAMPLE_BLOCK")
	require.Contains(t, updatedContent, "new content")
	require.Contains(t, updatedContent, "# END EXAMPLE_BLOCK")
	require.NotContains(t, updatedContent, "old content")
	require.True(t, contextutils.IsChanged(ctx))
}

func TestBlockModifiedChangeIndicatorIsSet(t *testing.T) {
	ctx := getCtx()
	ctx = contextutils.WithChangeIndicator(ctx)

	content := "# BEGIN EXAMPLE_BLOCK\nold content\n# END EXAMPLE_BLOCK\n"

	_, err := stringsutils.BlockInString(ctx, content, "EXAMPLE_BLOCK", "new content")
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctx))
}

func TestBlockNotModifiedChangeIndicatorNotSet(t *testing.T) {
	ctx := getCtx()
	ctx = contextutils.WithChangeIndicator(ctx)

	content := "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK\n"

	_, err := stringsutils.BlockInString(ctx, content, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctx))
}

func TestBlockAddedToContentWithTrailingNewline(t *testing.T) {
	ctx := getCtx()
	ctx = contextutils.WithChangeIndicator(ctx)

	content := "some existing content\n"

	updatedContent, err := stringsutils.BlockInString(ctx, content, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctx))
	require.Contains(t, updatedContent, "some existing content\n")
	require.Contains(t, updatedContent, "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK")
}

func TestBlockAddedToContentWithoutTrailingNewline(t *testing.T) {
	ctx := getCtx()
	ctx = contextutils.WithChangeIndicator(ctx)

	content := "some existing content"

	updatedContent, err := stringsutils.BlockInString(ctx, content, "EXAMPLE_BLOCK", "example content")
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctx))
	require.Contains(t, updatedContent, "some existing content\n")
	require.Contains(t, updatedContent, "# BEGIN EXAMPLE_BLOCK\nexample content\n# END EXAMPLE_BLOCK")
}

func TestBlockNameEmptyReturnsError(t *testing.T) {
	ctx := getCtx()

	_, err := stringsutils.BlockInString(ctx, "some content", "", "example content")
	require.Error(t, err)
}
