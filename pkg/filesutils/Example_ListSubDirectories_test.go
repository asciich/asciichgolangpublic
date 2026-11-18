package filesutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// A simple example how to list the direct sub directories of a local directory.
func Test_ExampleListSubDirectories(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Create a temporary directory with some content using plain golang:
	rootDir, err := os.MkdirTemp("", "demo-*")
	require.NoError(t, err)
	defer os.RemoveAll(rootDir)

	require.NoError(t, os.MkdirAll(filepath.Join(rootDir, "alpha", "nested"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(rootDir, "beta"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(rootDir, "file.txt"), []byte("hello"), 0o644))

	// List the sub directories.
	// Without any options only the direct sub directories are returned.
	// Files are not part of the result:
	subDirectories, err := ListSubDirectories(ctx, rootDir, &parameteroptions.ListDirectoryOptions{})
	require.NoError(t, err)
	require.Equal(
		t,
		[]string{
			filepath.Join(rootDir, "alpha"),
			filepath.Join(rootDir, "beta"),
		},
		subDirectories,
	)
}
