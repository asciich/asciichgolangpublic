package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBase(t *testing.T) {
	tests := []struct {
		command        []string
		expectedOutput string
	}{
		{[]string{"echo", "hello"}, "hello\n"},
		{[]string{"echo", "hello world"}, "hello world\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				fileBase := FileBase{}

				parent, err := fileBase.GetParentFileForBaseClass()
				assert.Nil(parent)
				assert.ErrorIs(err, ErrFileBaseParentNotSet)
				assert.ErrorIs(err, ErrTracedError)
			},
		)
	}
}

func TestFileBaseEnsureLineInFile_testcase1(t *testing.T) {
	const verbose bool = true

	tests := []struct {
		fileToTest File
	}{
		{MustGetLocalFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose))},
		{MustGetLocalCommandExecutorFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose), verbose)},
	}

	defer func() {
		for _, t := range tests {
			t.fileToTest.Delete(verbose)
		}
	}()

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const testContent string = "hello\nworld\n"
				tt.fileToTest.MustWriteString(testContent, verbose)

				assert.EqualValues(
					testContent,
					tt.fileToTest.MustReadAsString(),
				)

				tt.fileToTest.MustEnsureLineInFile("hello", verbose)
				assert.EqualValues(
					testContent,
					tt.fileToTest.MustReadAsString(),
				)
				tt.fileToTest.MustEnsureLineInFile("hello\n", verbose)
				assert.EqualValues(
					testContent,
					tt.fileToTest.MustReadAsString(),
				)
				tt.fileToTest.MustEnsureLineInFile("\nhello", verbose)
				assert.EqualValues(
					testContent,
					tt.fileToTest.MustReadAsString(),
				)
				tt.fileToTest.MustEnsureLineInFile("\nhello\n", verbose)
				assert.EqualValues(
					testContent,
					tt.fileToTest.MustReadAsString(),
				)
				tt.fileToTest.MustEnsureLineInFile("abc", verbose)
				assert.EqualValues(
					testContent+"abc\n",
					tt.fileToTest.MustReadAsString(),
				)
			},
		)
	}
}
