package files_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
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

				localPath, err := tempFile.GetLocalPath()
				require.NoError(t, err)
				err = files.Files().WriteStringToFile(localPath, tt.content, verbose)
				require.NoError(t, err)

				require.EqualValues(t, tempFile.MustReadAsString(), tt.content)
			},
		)
	}
}
