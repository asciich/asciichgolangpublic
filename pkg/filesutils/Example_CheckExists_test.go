package filesutils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

// Test_Example_CheckExists verifies that CheckExists returns an error when the
// path does not exist and succeeds otherwise.
func Test_Example_CheckExists(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// --- Existing file ---
	tmpFile, err := os.CreateTemp("", "filesutils_test_file_*")
	require.NoError(t, err)
	filePath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(filePath)

	err = filesutils.CheckExists(ctx, filePath)
	require.NoError(t, err, "CheckExists should succeed for existing file")

	// Remove the file and expect an error.
	os.Remove(filePath)
	err = filesutils.CheckExists(ctx, filePath)
	require.Error(t, err, "CheckExists should error for missing file")

	// --- Existing directory ---
	tempDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer filesutils.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

	err = filesutils.CheckExists(ctx, tempDir)
	require.NoError(t, err, "CheckExists should succeed for existing directory")

	// Delete directory and expect an error.
	err = os.RemoveAll(tempDir)
	require.NoError(t, err)
	err = filesutils.CheckExists(ctx, tempDir)
	require.Error(t, err, "CheckExists should error for missing directory")
}
