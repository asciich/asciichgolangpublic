package tarutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/archiveutils/tarutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const fileName = "file_name.txt"

				tarArchiveBytes, err := tarutils.CreateTarArchiveFromFileContentStringAndGetAsBytes(
					fileName,
					tt.content,
				)
				require.NoError(t, err)

				readContent, err := tarutils.ReadFileFromTarArchiveBytesAsString(
					tarArchiveBytes,
					fileName,
				)
				require.NoError(t, err)

				require.EqualValues(t, tt.content, readContent)
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const fileName = "file_name.txt"
				const fileName2 = "file_name2.txt"

				tarArchiveBytes, err := tarutils.CreateTarArchiveFromFileContentStringAndGetAsBytes(
					fileName,
					tt.content,
				)
				require.NoError(t, err)

				tarArchiveBytes, err = tarutils.AddFileFromFileContentStringToTarArchiveBytes(
					tarArchiveBytes,
					fileName2,
					tt.content+"2",
				)
				require.NoError(t, err)

				require.EqualValues(
					t,
					tt.content,
					mustutils.Must(tarutils.ReadFileFromTarArchiveBytesAsString(
						tarArchiveBytes,
						fileName,
					)),
				)
				require.EqualValues(
					t,
					tt.content+"2",
					mustutils.Must(tarutils.ReadFileFromTarArchiveBytesAsString(
						tarArchiveBytes,
						fileName2,
					)),
				)

				fileList, err := tarutils.ListFileNamesFromTarArchiveBytes(tarArchiveBytes)
				require.NoError(t, err)
				require.EqualValues(t, []string{"file_name.txt", "file_name2.txt"}, fileList)

			},
		)
	}
}
