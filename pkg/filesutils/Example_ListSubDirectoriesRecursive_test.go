package filesutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// An example how to list all sub directories recursively.
func Test_ExampleListSubDirectories_recursive(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Create a temporary directory with some nested content using plain golang:
	rootDir, err := os.MkdirTemp("", "demo-*")
	require.NoError(t, err)
	defer os.RemoveAll(rootDir)

	require.NoError(t, os.MkdirAll(filepath.Join(rootDir, "alpha", "nested"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(rootDir, "beta"), 0o755))

	// With the Recursive option all nested sub directories are returned as well.
	// The given root directory itself is never part of the result:
	subDirectories, err := ListSubDirectories(ctx, rootDir, &parameteroptions.ListDirectoryOptions{
		Recursive: true,
	})
	require.NoError(t, err)
	require.Equal(
		t,
		[]string{
			filepath.Join(rootDir, "alpha"),
			filepath.Join(rootDir, "alpha", "nested"),
			filepath.Join(rootDir, "beta"),
		},
		subDirectories,
	)
}
