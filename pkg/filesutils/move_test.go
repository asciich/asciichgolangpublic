package filesutils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestFile_MoveToPath tests MoveToPath method
func TestFile_MoveToPath(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativefilesoo"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				// Create source file
				srcFile := getTemporaryFileToTest(tt.implementationName)
				defer srcFile.Delete(ctx, &filesoptions.DeleteOptions{})

				// Write content to source
				err := srcFile.WriteString(ctx, "test content", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// Get source path
				srcPath, err := srcFile.GetPath()
				require.NoError(t, err)

				// Create destination path
				destPath := srcPath + ".moved"
				defer os.Remove(destPath)

				// Move file
				movedFile, err := srcFile.MoveToPath(ctx, destPath, false)
				require.NoError(t, err)
				require.NotNil(t, movedFile)

				// Verify source no longer exists
				exists, err := srcFile.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				// Verify destination exists
				exists, err = movedFile.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				// Verify content
				content, err := movedFile.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "test content", content)
			},
		)
	}
}
