package filesgeneric_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func createTempFileAndGetPath() (path string) {
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
func getFileToTest(implementationName string) (file filesinterfaces.File) {
	if implementationName == "localFile" {
		return mustutils.Must(files.GetLocalFileByPath(createTempFileAndGetPath()))
	}

	if implementationName == "localCommandExecutorFile" {
		return mustutils.Must(files.GetLocalCommandExecutorFileByPath(createTempFileAndGetPath()))
	}

	logging.LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	return nil
}

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

				fileBase := filesgeneric.FileBase{}

				parent, err := fileBase.GetParentFileForBaseClass()
				require.Nil(parent)
				require.ErrorIs(err, filesgeneric.ErrFileBaseParentNotSet)
				require.ErrorIs(err, tracederrors.ErrTracedError)
			},
		)
	}
}

func TestFileBaseEnsureLineInFile_testcase1(t *testing.T) {
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

				content, err := fileToTest.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, testContent, content)

				content, err = fileToTest.ReadAsString(ctx)
				require.NoError(t, err)
				err = fileToTest.EnsureLineInFile(ctx, "hello")
				require.NoError(t, err)
				require.EqualValues(t, testContent, content)

				content, err = fileToTest.ReadAsString(ctx)
				require.NoError(t, err)
				err = fileToTest.EnsureLineInFile(ctx, "hello\n")
				require.NoError(t, err)
				require.EqualValues(t, testContent, content)

				content, err = fileToTest.ReadAsString(ctx)
				require.NoError(t, err)
				err = fileToTest.EnsureLineInFile(ctx, "\nhello")
				require.NoError(t, err)
				require.EqualValues(t, testContent, content)

				content, err = fileToTest.ReadAsString(ctx)
				require.NoError(t, err)
				err = fileToTest.EnsureLineInFile(ctx, "\nhello\n")
				require.NoError(t, err)
				require.EqualValues(t, testContent, content)

				err = fileToTest.EnsureLineInFile(ctx, "abc")
				require.NoError(t, err)

				content, err = fileToTest.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, testContent+"abc\n", content)
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
					err = fileToTest.EnsureLineInFile(ctx, tt.line)
					require.NoError(t, err)

					content, err := fileToTest.ReadAsString(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.expected, content)
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
				err = toTest.RemoveLinesWithPrefix(ctx, tt.prefix)
				require.NoError(t, err)

				content, err := toTest.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedOutput, content)
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

				value, err := toTest.GetValueAsString(ctx, tt.key)
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

				intValue, err := toTest.GetValueAsInt(ctx, tt.key)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedValue, intValue)
			},
		)
	}
}
