package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestFilesWriteStringToFile(t *testing.T) {
	tests := []struct {
		implementationName string
		content string
	}{
		{"localFile", "hello_world"},
		{"commandExecutorFile", "hello_world"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				tempFile := getFileToTest(tt.implementationName)
				defer tempFile.Delete(verbose)

				tempFile2 := getFileToTest(tt.implementationName)
				defer tempFile2.Delete(verbose)

				Files().MustWriteStringToFile(tempFile.MustGetLocalPath(), tt.content, verbose)

				// Since used often there is a convenience function to write a file by path:
				content2 := tt.content + "2"
				MustWriteStringToFile(tempFile2.MustGetLocalPath(), content2, verbose)

				assert.EqualValues(
					tempFile.MustReadAsString(),
					tt.content,
				)

				assert.EqualValues(
					tempFile2.MustReadAsString(),
					content2,
				)

				assert.EqualValues(
					Files().MustReadAsString(tempFile.MustGetLocalPath()),
					tt.content,
				)

				// Since used often there is a convenience function to read a file by path:
				assert.EqualValues(
					MustReadFileAsString(tempFile.MustGetLocalPath()),
					tt.content,
				)
			},
		)
	}
}
