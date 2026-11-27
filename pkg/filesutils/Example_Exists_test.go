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

func Test_Example_Exists(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// --- Test with a regular file ---
	tmpFile, err := os.CreateTemp("", "filesutils_test_file_*")
	require.NoError(t, err)
	filePath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(filePath)

	exists, err := filesutils.Exists(ctx, filePath)
	require.NoError(t, err)
	require.True(t, exists, "%s should exist after creation", filePath)

	// Remove the file and ensure Exists reports false.
	os.Remove(filePath)
	exists, err = filesutils.Exists(ctx, filePath)
	require.NoError(t, err)
	require.False(t, exists, "%s should not exist after removal", filePath)

	// --- Test with a directory ---
	tempDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer filesutils.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

	exists, err = filesutils.Exists(ctx, tempDir)
	require.NoError(t, err)
	require.True(t, exists, "%s should exist as a directory", tempDir)

	// Delete the directory and confirm non‑existence.
	err = os.RemoveAll(tempDir)
	require.NoError(t, err)
	exists, err = filesutils.Exists(ctx, tempDir)
	require.NoError(t, err)
	require.False(t, exists, "%s should not exist after deletion", tempDir)
}
