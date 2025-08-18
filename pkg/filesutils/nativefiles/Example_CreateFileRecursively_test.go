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

func Test_Example_CreateFileRecursively(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.WithVerbose(context.TODO())

	// To not clutter the file system we run this example in a temporary directory:
	testRootDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	// After the test run the test directory should be deleted:
	defer nativefiles.Delete(ctx, testRootDir, &filesoptions.DeleteOptions{})

	// Define the path to the directory to create with multiple subdirs in between:
	testFile := filepath.Join(testRootDir, "a", "b", "c", "testfile.txt")

	// The testFile should not exist yet:
	require.False(t, nativefiles.Exists(ctx, testFile))

	// Create the file
	err = nativefiles.Create(ctx, testFile)
	require.NoError(t, err)

	// The testFile should now exists:
	require.True(t, nativefiles.Exists(ctx, testFile))

	// Just for visualization: The sub directories in between exist of course as well:
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir)))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a")))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a", "b")))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a", "b", "c")))
	require.True(t, nativefiles.Exists(ctx, filepath.Join(testRootDir, "a", "b", "c", "testfile.txt")))
}
