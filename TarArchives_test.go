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

func TestTarArchiveAddAndGetFileOnTarBytes_multipleFiles(t *testing.T) {

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
				const fileName2 = "file_name2.txt"

				tarArchiveBytes := TarArchives().MustCreateTarArchiveFromFileContentStringAndGetAsBytes(
					fileName,
					tt.content,
				)

				tarArchiveBytes = TarArchives().MustAddFileFromFileContentStringToTarArchiveBytes(
					tarArchiveBytes,
					fileName2,
					tt.content+"2",
				)

				assert.EqualValues(
					tt.content,
					TarArchives().MustReadFileFromTarArchiveBytesAsString(
						tarArchiveBytes,
						fileName,
					),
				)
				assert.EqualValues(
					tt.content+"2",
					TarArchives().MustReadFileFromTarArchiveBytesAsString(
						tarArchiveBytes,
						fileName2,
					),
				)
			},
		)
	}
}