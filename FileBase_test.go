package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func TestFileBase(t *testing.T) {
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

				fileBase := FileBase{}

				parent, err := fileBase.GetParentFileForBaseClass()
				assert.Nil(parent)
				assert.ErrorIs(err, ErrFileBaseParentNotSet)
				assert.ErrorIs(err, tracederrors.ErrTracedError)
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
			testutils.MustFormatAsTestname(tt),
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

func TestFileBase_EnsureLineInFile_testcaseWriteToNonexstingString(t *testing.T) {
	const verbose bool = true

	tests := []struct {
		fileToTest File
		line       string
		expected   string
	}{
		{
			MustGetLocalFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)),
			"hello",
			"hello\n",
		},
		{
			MustGetLocalFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose)),
			"hello world",
			"hello world\n",
		},
		{
			MustGetLocalCommandExecutorFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose), verbose),
			"hello",
			"hello\n",
		},
		{
			MustGetLocalCommandExecutorFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose), verbose),
			"hello world",
			"hello world\n",
		},
	}

	defer func() {
		for _, t := range tests {
			t.fileToTest.Delete(verbose)
		}
	}()

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				tt.fileToTest.MustDelete(verbose)

				for i := 0; i < 2; i++ {
					tt.fileToTest.MustEnsureLineInFile(tt.line, verbose)
					assert.EqualValues(
						tt.expected,
						tt.fileToTest.MustReadAsString(),
					)
				}
			},
		)
	}
}

func getTemporaryFileToTest(implementationName string) (file File) {
	const verbose = true

	if implementationName == "localFile" {
		return MustGetLocalFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose))
	} else if implementationName == "commandExecutorFile" {
		return MustGetLocalCommandExecutorFileByFile(TemporaryFiles().MustCreateEmptyTemporaryFile(verbose), verbose)
	}

	logging.LogFatalWithTracef("Unknown implementation name '%s'", implementationName)

	return nil
}

func TestFileBase_RemoveLinesWithPrefix(t *testing.T) {
	tests := []struct {
		implementationName string
		input              string
		prefix             string
		expectedOutput     string
	}{
		// Test LocalFile implementation
		{"localFile", "", "abc", ""},
		{"localFile", "\n", "abc", "\n"},
		{"localFile", "abc\n", "abc", ""},
		{"localFile", "1: a\n2: b\n3: c\n", "1", "2: b\n3: c\n"},
		{"localFile", "1: a\n2: b\n3: c", "1", "2: b\n3: c"},
		{"localFile", "1: a\n2: b\n3: c\n", "2", "1: a\n3: c\n"},
		{"localFile", "1: a\n2: b\n3: c", "2", "1: a\n3: c"},
		{"localFile", "1: a\n2: b\n3: c\n", "2:", "1: a\n3: c\n"},
		{"localFile", "1: a\n2: b\n3: c", "2:", "1: a\n3: c"},
		{"localFile", "1: a\n2: b\n3: c\n", "2: ", "1: a\n3: c\n"},
		{"localFile", "1: a\n2: b\n3: c", "2: ", "1: a\n3: c"},
		{"localFile", "1: a\n2: b\n3: c\n", "3", "1: a\n2: b\n"},
		{"localFile", "1: a\n2: b\n3: c", "3", "1: a\n2: b"},
		{"localFile", "1: a\n2: b\n3: c\n", "3:", "1: a\n2: b\n"},
		{"localFile", "1: a\n2: b\n3: c", "3:", "1: a\n2: b"},
		{"localFile", "1: a\n2: b\n3: c\n", "3: ", "1: a\n2: b\n"},
		{"localFile", "1: a\n2: b\n3: c", "3: ", "1: a\n2: b"},

		// Test CommandExecutorFile implementation
		{"commandExecutorFile", "", "abc", ""},
		{"commandExecutorFile", "\n", "abc", "\n"},
		{"commandExecutorFile", "abc\n", "abc", ""},
		{"commandExecutorFile", "1: a\n2: b\n3: c\n", "1", "2: b\n3: c\n"},
		{"commandExecutorFile", "1: a\n2: b\n3: c", "1", "2: b\n3: c"},
		{"commandExecutorFile", "1: a\n2: b\n3: c\n", "2", "1: a\n3: c\n"},
		{"commandExecutorFile", "1: a\n2: b\n3: c", "2", "1: a\n3: c"},
		{"commandExecutorFile", "1: a\n2: b\n3: c\n", "2:", "1: a\n3: c\n"},
		{"commandExecutorFile", "1: a\n2: b\n3: c", "2:", "1: a\n3: c"},
		{"commandExecutorFile", "1: a\n2: b\n3: c\n", "2: ", "1: a\n3: c\n"},
		{"commandExecutorFile", "1: a\n2: b\n3: c", "2: ", "1: a\n3: c"},
		{"commandExecutorFile", "1: a\n2: b\n3: c\n", "3", "1: a\n2: b\n"},
		{"commandExecutorFile", "1: a\n2: b\n3: c", "3", "1: a\n2: b"},
		{"commandExecutorFile", "1: a\n2: b\n3: c\n", "3:", "1: a\n2: b\n"},
		{"commandExecutorFile", "1: a\n2: b\n3: c", "3:", "1: a\n2: b"},
		{"commandExecutorFile", "1: a\n2: b\n3: c\n", "3: ", "1: a\n2: b\n"},
		{"commandExecutorFile", "1: a\n2: b\n3: c", "3: ", "1: a\n2: b"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				toTest := getTemporaryFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

				toTest.MustWriteString(tt.input, verbose)
				toTest.MustRemoveLinesWithPrefix(tt.prefix, verbose)

				assert.EqualValues(
					tt.expectedOutput,
					toTest.MustReadAsString(),
				)
			},
		)
	}
}

func TestFileBase_GetValueAsString(t *testing.T) {
	tests := []struct {
		implementationName string
		input              string
		key                string
		expectedValue      string
	}{
		{"localFile", "a=b\nc=hello world", "a", "b"},
		{"localFile", "a=c\nc=hello world", "c", "hello world"},
		{"commandExecutorFile", "a=b\nc=hello world", "a", "b"},
		{"commandExecutorFile", "a=c\nc=hello world", "c", "hello world"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				toTest := getTemporaryFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

				toTest.MustWriteString(tt.input, verbose)

				assert.EqualValues(
					tt.expectedValue,
					toTest.MustGetValueAsString(tt.key),
				)
			},
		)
	}
}

func TestFileBase_GetValueAsInt(t *testing.T) {
	tests := []struct {
		implementationName string
		input              string
		key                string
		expectedValue      int
	}{
		{"localFile", "a=1\nb=0\nc=-5\n", "a", 1},
		{"localFile", "a=1\nb=0\nc=-5\n", "b", 0},
		{"localFile", "a=1\nb=0\nc=-5\n", "c", -5},
		{"commandExecutorFile", "a=1\nb=0\nc=-5\n", "a", 1},
		{"commandExecutorFile", "a=1\nb=0\nc=-5\n", "b", 0},
		{"commandExecutorFile", "a=1\nb=0\nc=-5\n", "c", -5},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				toTest := getTemporaryFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

				toTest.MustWriteString(tt.input, verbose)

				assert.EqualValues(
					tt.expectedValue,
					toTest.MustGetValueAsInt(tt.key),
				)
			},
		)
	}
}
