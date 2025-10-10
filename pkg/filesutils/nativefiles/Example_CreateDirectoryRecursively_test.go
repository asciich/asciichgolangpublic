package nativefiles_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_Example_CreateDirectoryRecursively(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.WithVerbose(context.TODO())

	// To not clutter the file system we run this example in a temporary directory:
	testRootDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	// After the test run the test directory should be deleted:
	defer nativefiles.Delete(ctx, testRootDir, &filesoptions.DeleteOptions{})

	// Define the path to the directory to create with multiple subdirs in between:
	testDir := filepath.Join(testRootDir, "a", "b", "c", "testdir")

	// The testDir should not exist yet:
	require.False(t, nativefiles.Exists(ctx, testDir))

	// Create the directory
	err = nativefiles.CreateDirectory(ctx, testDir, &filesoptions.CreateOptions{})
	require.NoError(t, err)

	// The testDir should now exists:
	require.True(t, nativefiles.Exists(ctx, testDir))

	// Just for visualization: The sub directories in between exist of course as well:
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir)))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a")))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a", "b")))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a", "b", "c")))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a", "b", "c", "testdir")))
}
