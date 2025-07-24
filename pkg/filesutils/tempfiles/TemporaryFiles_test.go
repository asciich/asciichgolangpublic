package tempfiles

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestCreateTemporaryFile(t *testing.T) {
	tests := []struct {
		content string
	}{
		{""},
		{"a"},
		{"hello world"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				file := MustCreateFromString(tt.content, verbose)
				defer file.Delete(verbose)

				require.True(file.MustExists(verbose))
				require.EqualValues(tt.content, file.MustReadAsString())
			},
		)
	}
}

func TestCreateEmptyTemporaryFile(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				file := MustCreateEmptyTemporaryFile(verbose)
				defer file.Delete(verbose)

				require.True(file.MustExists(verbose))
				require.EqualValues("", file.MustReadAsString())
				require.True(strings.HasPrefix(file.MustGetLocalPath(), "/tmp/"))
			},
		)
	}
}

func TestCreateEmptyTemporaryFileAndGetPath(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				filePath := MustCreateEmptyTemporaryFileAndGetPath(verbose)
				file := files.MustNewLocalFileByPath(filePath)
				defer file.Delete(verbose)

				require.True(file.MustExists(verbose))
				require.EqualValues("", file.MustReadAsString())
				require.True(strings.HasPrefix(filePath, "/tmp/"))
				require.True(strings.HasPrefix(file.MustGetPath(), "/tmp/"))
			},
		)
	}
}

func getFileToTest(implementationName string) (fileToTest files.File) {
	temporayFile := MustCreateEmptyTemporaryFileAndGetPath(false)

	if implementationName == "localFile" {
		return files.MustGetLocalFileByPath(temporayFile)
	}

	if implementationName == "localCommandExecutorFile" {
		return files.MustGetLocalCommandExecutorFileByPath(temporayFile)
	}

	logging.LogFatalWithTracef("Unknown implementation name '%s'", implementationName)
	return nil
}

func TestTemporaryFilesCreateFromFile(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
	}{
		{"localFile", "testcase"},
		{"localCommandExecutorFile", "testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				sourceFile := getFileToTest(tt.implementationName)
				sourceFile.MustWriteString(tt.content, verbose)
				defer sourceFile.Delete(verbose)

				require.EqualValues(
					tt.content,
					sourceFile.MustReadAsString(),
				)

				tempFile := MustCreateTemporaryFileFromFile(sourceFile, verbose)
				defer tempFile.Delete(verbose)

				require.EqualValues(
					tt.content,
					tempFile.MustReadAsString(),
				)
			},
		)
	}
}
