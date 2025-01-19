package tempfiles

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				assert := assert.New(t)

				const verbose = true

				file := MustCreateFromString(tt.content, verbose)
				defer file.Delete(verbose)

				assert.True(file.MustExists(verbose))
				assert.EqualValues(tt.content, file.MustReadAsString())
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
				assert := assert.New(t)

				const verbose bool = true

				file := MustCreateEmptyTemporaryFile(verbose)
				defer file.Delete(verbose)

				assert.True(file.MustExists(verbose))
				assert.EqualValues("", file.MustReadAsString())
				assert.True(strings.HasPrefix(file.MustGetLocalPath(), "/tmp/"))
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
				assert := assert.New(t)

				const verbose bool = true

				filePath := MustCreateEmptyTemporaryFileAndGetPath(verbose)
				file := files.MustNewLocalFileByPath(filePath)
				defer file.Delete(verbose)

				assert.True(file.MustExists(verbose))
				assert.EqualValues("", file.MustReadAsString())
				assert.True(strings.HasPrefix(filePath, "/tmp/"))
				assert.True(strings.HasPrefix(file.MustGetPath(), "/tmp/"))
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
				assert := assert.New(t)

				const verbose bool = true

				sourceFile := getFileToTest(tt.implementationName)
				sourceFile.MustWriteString(tt.content, verbose)
				defer sourceFile.Delete(verbose)

				assert.EqualValues(
					tt.content,
					sourceFile.MustReadAsString(),
				)

				tempFile := MustCreateTemporaryFileFromFile(sourceFile, verbose)
				defer tempFile.Delete(verbose)

				assert.EqualValues(
					tt.content,
					tempFile.MustReadAsString(),
				)
			},
		)
	}
}
