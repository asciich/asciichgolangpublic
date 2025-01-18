package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.MustDelete(verbose)

				const testContent string = "hello\nworld\n"
				fileToTest.MustWriteString(testContent, verbose)

				assert.EqualValues(
					testContent,
					fileToTest.MustReadAsString(),
				)

				fileToTest.MustEnsureLineInFile("hello", verbose)
				assert.EqualValues(
					testContent,
					fileToTest.MustReadAsString(),
				)
				fileToTest.MustEnsureLineInFile("hello\n", verbose)
				assert.EqualValues(
					testContent,
					fileToTest.MustReadAsString(),
				)
				fileToTest.MustEnsureLineInFile("\nhello", verbose)
				assert.EqualValues(
					testContent,
					fileToTest.MustReadAsString(),
				)
				fileToTest.MustEnsureLineInFile("\nhello\n", verbose)
				assert.EqualValues(
					testContent,
					fileToTest.MustReadAsString(),
				)
				fileToTest.MustEnsureLineInFile("abc", verbose)
				assert.EqualValues(
					testContent+"abc\n",
					fileToTest.MustReadAsString(),
				)
			},
		)
	}
}

func TestFileBase_EnsureLineInFile_testcaseWriteToNonexstingString(t *testing.T) {
	const verbose bool = true

	tests := []struct {
		implementationName string
		line               string
		expected           string
	}{
		{
			"localFile",
			"hello",
			"hello\n",
		},
		{
			"localFile",
			"hello world",
			"hello world\n",
		},
		{
			"localFile",
			"hello",
			"hello\n",
		},
		{
			"localFile",
			"hello world",
			"hello world\n",
		},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.MustDelete(true)

				fileToTest.MustDelete(verbose)

				for i := 0; i < 2; i++ {
					fileToTest.MustEnsureLineInFile(tt.line, verbose)
					assert.EqualValues(
						tt.expected,
						fileToTest.MustReadAsString(),
					)
				}
			},
		)
	}
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
		{"localCommandExecutorFile", "", "abc", ""},
		{"localCommandExecutorFile", "\n", "abc", "\n"},
		{"localCommandExecutorFile", "abc\n", "abc", ""},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c\n", "1", "2: b\n3: c\n"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c", "1", "2: b\n3: c"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c\n", "2", "1: a\n3: c\n"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c", "2", "1: a\n3: c"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c\n", "2:", "1: a\n3: c\n"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c", "2:", "1: a\n3: c"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c\n", "2: ", "1: a\n3: c\n"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c", "2: ", "1: a\n3: c"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c\n", "3", "1: a\n2: b\n"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c", "3", "1: a\n2: b"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c\n", "3:", "1: a\n2: b\n"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c", "3:", "1: a\n2: b"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c\n", "3: ", "1: a\n2: b\n"},
		{"localCommandExecutorFile", "1: a\n2: b\n3: c", "3: ", "1: a\n2: b"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				toTest := getFileToTest(tt.implementationName)
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
		{"localCommandExecutorFile", "a=b\nc=hello world", "a", "b"},
		{"localCommandExecutorFile", "a=c\nc=hello world", "c", "hello world"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				toTest := getFileToTest(tt.implementationName)
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
		{"localCommandExecutorFile", "a=1\nb=0\nc=-5\n", "a", 1},
		{"localCommandExecutorFile", "a=1\nb=0\nc=-5\n", "b", 0},
		{"localCommandExecutorFile", "a=1\nb=0\nc=-5\n", "c", -5},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				toTest := getFileToTest(tt.implementationName)
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
