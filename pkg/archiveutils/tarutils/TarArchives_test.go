package tarutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

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

func Test_GetFileFromArchive(t *testing.T) {
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
				ctx := getCtx()

				const fileName = "file_name.txt"

				tarArchiveBytes, err := tarutils.CreateTarArchiveFromFileContentStringAndGetAsBytes(
					fileName,
					tt.content,
				)
				require.NoError(t, err)

				archivePath, err := tempfiles.CreateTemporaryFileFromContentBytes(ctx, tarArchiveBytes)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, archivePath, &filesoptions.DeleteOptions{})

				readContent, err := tarutils.ReadFileFromTarArchiveAsBytes(
					ctx,
					archivePath,
					fileName,
				)
				require.NoError(t, err)

				require.EqualValues(t, tt.content, readContent)
			},
		)
	}
}

func Test_ExtractFileFromArchive(t *testing.T) {
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
				ctx := getCtx()

				const fileName = "file_name.txt"

				tarArchiveBytes, err := tarutils.CreateTarArchiveFromFileContentStringAndGetAsBytes(
					fileName,
					tt.content,
				)
				require.NoError(t, err)

				archivePath, err := tempfiles.CreateTemporaryFileFromContentBytes(ctx, tarArchiveBytes)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, archivePath, &filesoptions.DeleteOptions{})

				destPath, err := tempfiles.CreateTemporaryFile(ctx)
				require.NoError(t,err)
				defer nativefiles.Delete(ctx, destPath, &filesoptions.DeleteOptions{})

				err = tarutils.ExtractFileFromTarArchive(
					ctx,
					archivePath,
					fileName,
					destPath,
				)
				require.NoError(t, err)

				readContent, err := nativefiles.ReadAsBytes(ctx, destPath)
				require.NoError(t, err)
				require.EqualValues(t, tt.content, readContent)
			},
		)
	}
}

func Test_GetFileFromTarGzArchive(t *testing.T) {
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
				ctx := getCtx()

				const fileName = "file_name.txt"

				tarGzArchiveBytes, err := tarutils.CreateTarGzArchiveFromFileContentStringAndGetAsBytes(
					fileName,
					tt.content,
				)
				require.NoError(t, err)

				archivePath, err := tempfiles.CreateTemporaryFileFromContentBytes(ctx, tarGzArchiveBytes)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, archivePath, &filesoptions.DeleteOptions{})

				readContent, err := tarutils.ReadFileFromTarArchiveAsBytes(
					ctx,
					archivePath,
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
