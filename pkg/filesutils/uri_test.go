package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestFile_GetUriAsString tests GetUriAsString method
func TestFile_GetUriAsString(t *testing.T) {
	tests := []struct {
		implementationName string
		filePath           string
	}{
		{"nativefilesoo", "/tmp/testfile.txt"},
		{"commandExecutorFileExec", "/tmp/testfile.txt"},
		{"commandExecutorFileBash", "/tmp/testfile.txt"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName, tt.filePath)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				uri, err := fileToTest.GetUriAsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.filePath, uri)
			},
		)
	}
}

// TestFile_String tests String method
func TestFile_String(t *testing.T) {
	tests := []struct {
		implementationName string
		filePath           string
	}{
		{"nativefilesoo", "/tmp/testfile.txt"},
		{"commandExecutorFileExec", "/tmp/testfile.txt"},
		{"commandExecutorFileBash", "/tmp/testfile.txt"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				fileToTest := getFileToTest(tt.implementationName, tt.filePath)
				defer fileToTest.Delete(getCtx(), &filesoptions.DeleteOptions{})

				path := fileToTest.String()
				require.EqualValues(t, tt.filePath, path)
			},
		)
	}
}
