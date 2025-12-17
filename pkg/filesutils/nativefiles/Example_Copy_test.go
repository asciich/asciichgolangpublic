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

// This Examples shows how to use the nativefiles.Copy function.
func Test_Example_copy(t *testing.T) {
	// Use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Let's create a simple source file:
	src, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
	require.NoError(t, err)
	// Delete the source file when the test is finished:
	defer nativefiles.Delete(ctx, src, &filesoptions.DeleteOptions{})

	// Let's create a temprary directory where we can copy the src file to:
	dstDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	// Delet the dstDir when the test is finished:
	defer nativefiles.Delete(ctx, dstDir, &filesoptions.DeleteOptions{})
	
	// Define the destination path:
	dst := filepath.Join(dstDir, "dst")
	// The dst file does not exist yet:
	require.False(t, nativefiles.Exists(ctx, dst))

	// Copy the file:
	err = nativefiles.Copy(ctx, src, dst, &filesoptions.CopyOptions{})
	require.NoError(t, err)

	// Now both files exist:
	require.True(t, nativefiles.Exists(ctx, src))
	require.True(t, nativefiles.Exists(ctx, dst))

	// And both file have the same content:
	content1, err := nativefiles.ReadAsString(ctx, src)
	require.NoError(t, err)
	content2, err := nativefiles.ReadAsString(ctx, dst)
	require.NoError(t, err)
	require.EqualValues(t, content1, content2)
	require.EqualValues(t, "hello world", content2)
}
