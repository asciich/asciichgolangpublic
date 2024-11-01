package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTarArchiveAddAndGetFileOnTarBytes(t *testing.T) {

	tests := []struct {
		content string
	}{
		{"hello"},
		{"multi\nline"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const fileName = "file_name.txt"

				tarArchiveBytes := TarArchives().MustCreateTarArchiveFromFileContentStringAndGetAsBytes(
					fileName,
					tt.content,
				)

				readContent := TarArchives().MustReadFileFromTarArchiveBytesAsString(
					tarArchiveBytes,
					fileName,
				)

				assert.EqualValues(
					tt.content,
					readContent,
				)
			},
		)
	}
}
