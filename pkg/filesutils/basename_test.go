package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestFile_GetBaseName tests GetBaseName method
func TestFile_GetBaseName(t *testing.T) {
	tests := []struct {
		implementationName string
		filePath           string
		expectedBaseName   string
	}{
		{"nativefilesoo", "/tmp/testfile.txt", "testfile.txt"},
		{"nativefilesoo", "/tmp/subdir/file.txt", "file.txt"},
		{"commandExecutorFileExec", "/tmp/testfile.txt", "testfile.txt"},
		{"commandExecutorFileBash", "/tmp/testfile.txt", "testfile.txt"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName, tt.filePath)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				baseName, err := fileToTest.GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedBaseName, baseName)
			},
		)
	}
}
