package nativefiles_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

// This example shows how to use the nativefiles.Move function.
func Test_Example_move(t *testing.T) {
	// Use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Let's create a simple source file:
	src, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
	require.NoError(t, err)
	// Delete the source file when the test is finished:
	defer nativefiles.Delete(ctx, src, &filesoptions.DeleteOptions{})

	// Let's create a temporary directory where we can move the src file to:
	dstDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	// Delete the dstDir when the test is finished:
	defer nativefiles.Delete(ctx, dstDir, &filesoptions.DeleteOptions{})

	// Define the destination path:
	dst := filepath.Join(dstDir, "dst")
	// The dst file does not exist yet:
	dstExists, err := nativefiles.Exists(ctx, dst)
	require.NoError(t, err)
	require.False(t, dstExists)

	// Move the file:
	err = nativefiles.Move(ctx, src, dst, &filesoptions.MoveOptions{})
	require.NoError(t, err)

	// Now only the dst file exists:
	srcExists, err := nativefiles.Exists(ctx, src)
	require.NoError(t, err)
	require.False(t, srcExists)

	dstExists, err = nativefiles.Exists(ctx, dst)
	require.NoError(t, err)
	require.True(t, dstExists)

	// And the dst file has the original content:
	content2, err := nativefiles.ReadAsString(ctx, dst, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.EqualValues(t, "hello world", content2)
}
