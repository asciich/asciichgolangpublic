package files

import (
	"testing"

	"github.com/stretchr/testify/require"
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
				require := require.New(t)

				const verbose bool = true

				tempFile := getFileToTest(tt.implementationName)
				defer tempFile.Delete(verbose)

				tempFile2 := getFileToTest(tt.implementationName)
				defer tempFile2.Delete(verbose)

				Files().MustWriteStringToFile(tempFile.MustGetLocalPath(), tt.content, verbose)

				// Since used often there is a convenience function to write a file by path:
				content2 := tt.content + "2"
				MustWriteStringToFile(tempFile2.MustGetLocalPath(), content2, verbose)

				require.EqualValues(
					tempFile.MustReadAsString(),
					tt.content,
				)

				require.EqualValues(
					tempFile2.MustReadAsString(),
					content2,
				)

				require.EqualValues(
					Files().MustReadAsString(tempFile.MustGetLocalPath()),
					tt.content,
				)

				// Since used often there is a convenience function to read a file by path:
				require.EqualValues(
					MustReadFileAsString(tempFile.MustGetLocalPath()),
					tt.content,
				)
			},
		)
	}
}
