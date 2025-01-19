package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func createTemFileAndGetPath() (path string) {
	file, err := os.CreateTemp("", "file_for_testing")
	if err != nil {
		logging.LogGoErrorFatalWithTrace(err)
	}

	return file.Name()
}

// Return a temporary file of the given 'implementationName'.
//
// Use defer file.Delete(verbose) to after calling this function to ensure
// the file is deleted after the test is over.
func getFileToTest(implementationName string) (file File) {
	if implementationName == "localFile" {
		return MustGetLocalFileByPath(createTemFileAndGetPath())
	}

	if implementationName == "localCommandExecutorFile" {
		return MustGetLocalCommandExecutorFileByPath(createTemFileAndGetPath())
	}

	logging.LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	return nil
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
			testutils.MustFormatAsTestname(tt),
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
			testutils.MustFormatAsTestname(tt),
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
			testutils.MustFormatAsTestname(tt),
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
			testutils.MustFormatAsTestname(tt),
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

func TestFile_MoveToPath(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localFile"},
		{"localCommandExecutorFile"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(verbose)
				destFile.Delete(verbose)

				assert.True(fileToTest.MustExists(verbose))
				assert.False(destFile.MustExists(verbose))

				movedFile := fileToTest.MustMoveToPath(destFile.MustGetPath(), false, verbose)

				assert.EqualValues(
					movedFile.MustGetPath(),
					destFile.MustGetPath(),
				)

				assert.EqualValues(
					movedFile.MustGetHostDescription(),
					destFile.MustGetHostDescription(),
				)

				assert.False(fileToTest.MustExists(verbose))
				assert.True(destFile.MustExists(verbose))
			},
		)
	}
}

func TestFile_CopyToFile(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
	}{
		{"localFile", "test content\nwith a new line\n"},
		{"localCommandExecutorFile", "test content\nwith a new line\n"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				srcFile := getFileToTest(tt.implementationName)
				srcFile.MustWriteString(tt.content, verbose)
				defer srcFile.Delete(verbose)

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(verbose)
				destFile.Delete(verbose)

				assert.True(srcFile.MustExists(verbose))
				assert.False(destFile.MustExists(verbose))

				srcFile.MustCopyToFile(destFile, verbose)

				assert.True(srcFile.MustExists(verbose))
				assert.True(destFile.MustExists(verbose))

				assert.EqualValues(
					tt.content,
					srcFile.MustReadAsString(),
				)

				assert.EqualValues(
					tt.content,
					destFile.MustReadAsString(),
				)
			},
		)
	}
}
