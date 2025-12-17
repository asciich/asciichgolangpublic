package nativefiles_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

// This Examples shows how to use the nativefiles.Cop function to copy a file as root using sudo.
func Test_Example_moveSudo(t *testing.T) {
	// Use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Let's create a simple source file:
	src, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
	require.NoError(t, err)
	// Delete the source file when the test is finished:
	defer nativefiles.Delete(ctx, src, &filesoptions.DeleteOptions{})

	// Define the destination path. /opt is usually not writeable as root and requires sudo.:
	dst := "/opt/tmp_dst_file_for_mv"
	// The dst file does not exist yet:
	require.False(t, nativefiles.Exists(ctx, dst))
	// Ensure the dst file is deleted after the test:
	defer nativefiles.Delete(ctx, dst, &filesoptions.DeleteOptions{UseSudo: true})

	// Moving the file without using sudo fails
	err = nativefiles.Move(ctx, src, dst, &filesoptions.MoveOptions{})
	require.Error(t, err)

	// But using sudo works:
	err = nativefiles.Move(ctx, src, dst, &filesoptions.MoveOptions{UseSudo: true})
	require.NoError(t, err)

	// Now only the dst files exist:
	require.False(t, nativefiles.Exists(ctx, src))
	require.True(t, nativefiles.Exists(ctx, dst))

	// And the dst file has the original content:
	content2, err := nativefiles.ReadAsString(ctx, dst)
	require.NoError(t, err)
	require.EqualValues(t, "hello world", content2)
}
