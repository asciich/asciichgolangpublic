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
	exists, err := nativefiles.Exists(ctx, testDir)
	require.NoError(t, err)
	require.False(t, exists)

	// Create the directory
	err = nativefiles.CreateDirectory(ctx, testDir, &filesoptions.CreateOptions{})
	require.NoError(t, err)

	// The testDir should now exist:
	exists, err = nativefiles.Exists(ctx, testDir)
	require.NoError(t, err)
	require.True(t, exists)

	// Just for visualization: The sub directories in between exist of course as well:
	for _, path := range []string{
		testRootDir,
		filepath.Join(testRootDir, "a"),
		filepath.Join(testRootDir, "a", "b"),
		filepath.Join(testRootDir, "a", "b", "c"),
		filepath.Join(testRootDir, "a", "b", "c", "testdir"),
	} {
		exists, err = nativefiles.Exists(ctx, path)
		require.NoError(t, err)
		require.True(t, exists, "'%s' must exist.", path)
	}
}
