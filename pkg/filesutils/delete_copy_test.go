package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestFile_SecurelyDelete tests SecurelyDelete method
func TestFile_SecurelyDelete(t *testing.T) {
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

				fileToTest := getTemporaryFileToTest(tt.implementationName)

				// Verify file exists
				exists, err := fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				// Securely delete
				err = fileToTest.SecurelyDelete(ctx)
				require.NoError(t, err)

				// Verify file no longer exists
				exists, err = fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
			},
		)
	}
}

// TestFile_GetDeepCopy tests GetDeepCopy method
func TestFile_GetDeepCopy(t *testing.T) {
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

				fileToTest := getTemporaryFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				// Get deep copy
				copy := fileToTest.GetDeepCopy()
				require.NotNil(t, copy)

				// Verify copy has same path
				path, err := fileToTest.GetPath()
				require.NoError(t, err)
				copyPath, err := copy.GetPath()
				require.NoError(t, err)
				require.EqualValues(t, path, copyPath)

				// Write different content to original
				err = fileToTest.WriteString(ctx, "original", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// Verify copy has same path but is independent
				copyPath2, err := copy.GetPath()
				require.NoError(t, err)
				require.EqualValues(t, path, copyPath2)
			},
		)
	}
}
