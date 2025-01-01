package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Return a temporary file of the given 'implementationName'.
//
// Use defer file.Delete(verbose) to after calling this function to ensure
// the file is deleted after the test is over.
func getFileToTest(implementationName string) (file File) {
	const verbose = true

	if implementationName == "localFile" {
		file = MustGetLocalFileByPath(
			TemporaryFiles().MustCreateEmptyTemporaryFileAndGetPath(verbose),
		)
	} else if implementationName == "localCommandExecutorFile" {
		file = MustGetLocalCommandExecutorFileByPath(
			TemporaryFiles().MustCreateEmptyTemporaryFileAndGetPath(verbose),
		)
	} else {
		LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	}

	return file
}

func TestFile_WriteString_ReadAsString(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
	}{
		{"localFile", "hello world"},
		{"localCommandExecutorFile", "hello world"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				fileToTest.MustWriteString(tt.content, verbose)

				assert.EqualValues(
					tt.content,
					fileToTest.MustReadAsString(),
				)
			},
		)
	}
}

func TestFile_Exists(t *testing.T) {

	tests := []struct {
		implementationName string
	}{
		{"localFile"},
		{"localCommandExecutorFile"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				assert.True(fileToTest.MustExists(verbose))

				fileToTest.MustDelete(verbose)

				assert.False(fileToTest.MustExists(verbose))
			},
		)
	}
}

func TestFile_Truncate(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localFile"},
		{"localCommandExecutorFile"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				for i := 0; i < 10; i++ {
					fileToTest.MustTruncate(int64(i), verbose)
					assert.EqualValues(
						fileToTest.MustGetSizeBytes(),
						int64(i),
					)
				}

				fileToTest.MustTruncate(0, verbose)

				assert.EqualValues(
					fileToTest.MustGetSizeBytes(),
					0,
				)
			},
		)
	}
}

func TestFile_ContainsLine(t *testing.T) {
	tests := []struct {
		implementationName string
		line               string
		expectedContains   bool
	}{
		{"localFile", "hello", false},
		{"localFile", "hello world", true},
		{"localCommandExecutorFile", "hello", false},
		{"localCommandExecutorFile", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				fileToTest.MustWriteString(
					"this is a\nhello world\nexample text.\n",
					verbose,
				)

				assert.EqualValues(
					tt.expectedContains,
					fileToTest.MustContainsLine(
						tt.line,
					),
				)
			},
		)
	}
}
