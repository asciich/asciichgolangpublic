package files_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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
				require := require.New(t)

				fileBase := files.FileBase{}

				parent, err := fileBase.GetParentFileForBaseClass()
				require.Nil(parent)
				require.ErrorIs(err, files.ErrFileBaseParentNotSet)
				require.ErrorIs(err, tracederrors.ErrTracedError)
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
				ctx := getCtx()
				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				const testContent string = "hello\nworld\n"
				err := fileToTest.WriteString(ctx, testContent, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.EqualValues(t, testContent, fileToTest.MustReadAsString())

				err = fileToTest.EnsureLineInFile("hello", verbose)
				require.NoError(t, err)
				require.EqualValues(t, testContent, fileToTest.MustReadAsString())

				err = fileToTest.EnsureLineInFile("hello\n", verbose)
				require.NoError(t, err)
				require.EqualValues(t, testContent, fileToTest.MustReadAsString())

				err = fileToTest.EnsureLineInFile("\nhello", verbose)
				require.NoError(t, err)
				require.EqualValues(t, testContent, fileToTest.MustReadAsString())

				err = fileToTest.EnsureLineInFile("\nhello\n", verbose)
				require.NoError(t, err)
				require.EqualValues(t, testContent, fileToTest.MustReadAsString())

				err = fileToTest.EnsureLineInFile("abc", verbose)
				require.NoError(t, err)
				require.EqualValues(t, testContent+"abc\n", fileToTest.MustReadAsString())
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
				ctx := getCtx()
				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				err := fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = fileToTest.EnsureLineInFile(tt.line, verbose)
					require.NoError(t, err)

					require.EqualValues(t, tt.expected, fileToTest.MustReadAsString())
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
				const verbose = true
				ctx := getCtx()

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(ctx, &filesoptions.DeleteOptions{})

				err := toTest.WriteString(ctx, tt.input, &filesoptions.WriteOptions{})
				require.NoError(t, err)
				toTest.MustRemoveLinesWithPrefix(tt.prefix, verbose)

				require.EqualValues(t, tt.expectedOutput, toTest.MustReadAsString())
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
				ctx := getCtx()

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(ctx, &filesoptions.DeleteOptions{})

				err := toTest.WriteString(ctx, tt.input, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				value, err := toTest.GetValueAsString(tt.key)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedValue, value)
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
				ctx := getCtx()

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(ctx, &filesoptions.DeleteOptions{})

				err := toTest.WriteString(ctx, tt.input, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				intValue, err := toTest.GetValueAsInt(tt.key)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedValue, intValue)
			},
		)
	}
}
