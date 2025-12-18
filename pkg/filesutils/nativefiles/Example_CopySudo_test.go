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
func Test_Example_copySudo(t *testing.T) {
	// Use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Let's create a simple source file:
	src, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
	require.NoError(t, err)
	// Delete the source file when the test is finished:
	defer nativefiles.Delete(ctx, src, &filesoptions.DeleteOptions{})

	// Define the destination path. /opt is usually not writeable as root and requires sudo.:
	dst := "/opt/tmp_dst_file"
	// The dst file does not exist yet:
	require.False(t, nativefiles.Exists(ctx, dst))
	// Ensure the dst file is deleted after the test:
	defer nativefiles.Delete(ctx, dst, &filesoptions.DeleteOptions{UseSudo: true})

	// Use sudo to copy to /opt/ which is not writeable as user:
	err = nativefiles.Copy(ctx, src, dst, &filesoptions.CopyOptions{UseSudo: true})
	require.NoError(t, err)

	// Now both files exist:
	require.True(t, nativefiles.Exists(ctx, src))
	require.True(t, nativefiles.Exists(ctx, dst))

	// And both file have the same content:
	content1, err := nativefiles.ReadAsString(ctx, src, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	content2, err := nativefiles.ReadAsString(ctx, dst,&filesoptions.ReadOptions{UseSudo: true})
	require.NoError(t, err)
	require.EqualValues(t, content1, content2)
	require.EqualValues(t, "hello world", content2)
}
