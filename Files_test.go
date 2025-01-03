package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesWriteStringToFile(t *testing.T) {
	tests := []struct {
		content string
	}{
		{"hello_world"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				tempFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
				defer tempFile.Delete(verbose)

				tempFile2 := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)
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

