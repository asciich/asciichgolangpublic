package filesutils_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestFile_GetParentDirectory tests GetParentDirectory method
func TestFile_GetParentDirectory(t *testing.T) {
	tests := []struct {
		implementationName string
		filePath           string
		expectedParentPath string
	}{
		{"nativefilesoo", "/tmp/testfile.txt", "/tmp"},
		{"nativefilesoo", "/tmp/subdir/file.txt", "/tmp/subdir"},
		{"commandExecutorFileExec", "/tmp/testfile.txt", "/tmp"},
		{"commandExecutorFileBash", "/tmp/testfile.txt", "/tmp"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName, tt.filePath)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				parentDir, err := fileToTest.GetParentDirectory(ctx)
				require.NoError(t, err)

				parentPath, err := parentDir.GetPath()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedParentPath, parentPath)
			},
		)
	}
}

// TestFile_GetLocalPathOrEmptyStringIfUnset tests GetLocalPathOrEmptyStringIfUnset method
func TestFile_GetLocalPathOrEmptyStringIfUnset(t *testing.T) {
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

				localPath, err := fileToTest.GetLocalPathOrEmptyStringIfUnset()
				require.NoError(t, err)
				require.NotEmpty(t, localPath)
				require.True(t, filepath.IsAbs(localPath))
			},
		)
	}
}
