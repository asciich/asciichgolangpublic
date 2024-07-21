package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesWriteStringToFile(t *testing.T) {
	tests := []struct {
		content string
	}{
		{"guguseli"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				tempFile := TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)

				Files().MustWriteStringToFile(tempFile.MustGetLocalPath(), tt.content, verbose)

				assert.EqualValues(
					tempFile.MustReadAsString(),
					tt.content,
				)
			},
		)
	}
}
