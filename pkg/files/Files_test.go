package files_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestFilesWriteStringToFile(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
	}{
		{"localFile", "hello_world"},
		{"localCommandExecutorFile", "hello_world"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				tempFile := getFileToTest(tt.implementationName)
				defer tempFile.Delete(verbose)

				tempFile2 := getFileToTest(tt.implementationName)
				defer tempFile2.Delete(verbose)

				err := files.Files().WriteStringToFile(tempFile.MustGetLocalPath(), tt.content, verbose)
				require.NoError(t, err)

				// Since used often there is a convenience function to write a file by path:
				content2 := tt.content + "2"
				err = files.WriteStringToFile(tempFile2.MustGetLocalPath(), content2, verbose)
				require.NoError(t, err)

				require.EqualValues(t, tempFile.MustReadAsString(), tt.content)

				require.EqualValues(t, tempFile2.MustReadAsString(), content2)

				require.EqualValues(t, mustutils.Must(files.Files().ReadAsString(tempFile.MustGetLocalPath())), tt.content)

				// Since used often there is a convenience function to read a file by path:
				require.EqualValues(t, mustutils.Must(files.ReadFileAsString(tempFile.MustGetLocalPath())), tt.content)
			},
		)
	}
}
