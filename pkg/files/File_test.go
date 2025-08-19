package files_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

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
		return files.MustGetLocalCommandExecutorFileByPath(createTempFileAndGetPath())
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
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				err := fileToTest.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.EqualValues(t, tt.content, fileToTest.MustReadAsString())
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
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				exists, err := fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				err = fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				exists, err = fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
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
				const verbose bool = true
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				for i := 0; i < 10; i++ {
					err := fileToTest.Truncate(int64(i), verbose)
					require.NoError(t, err)

					sizeBytes, err := fileToTest.GetSizeBytes()
					require.NoError(t, err)
					require.EqualValues(t, sizeBytes, int64(i))
				}

				err := fileToTest.Truncate(0, verbose)
				require.NoError(t, err)

				sizeBytes, err := fileToTest.GetSizeBytes()
				require.NoError(t, err)
				require.EqualValues(t, sizeBytes, 0)
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
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				err := fileToTest.WriteString(ctx, "this is a\nhello world\nexample text.\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				containsLine, err := fileToTest.ContainsLine(tt.line)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedContains, containsLine)
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
				const verbose bool = true
				ctx := getCtx()

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(ctx, &filesoptions.DeleteOptions{})
				destFile.Delete(ctx, &filesoptions.DeleteOptions{})

				exists, err := fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				exists, err = destFile.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				destFilePath, err := destFile.GetPath()
				require.NoError(t, err)

				movedFile, err := fileToTest.MoveToPath(destFilePath, false, verbose)
				require.NoError(t, err)

				movedFilePath, err := movedFile.GetPath()
				require.NoError(t, err)
				destFilePath, err = destFile.GetPath()
				require.NoError(t, err)

				require.EqualValues(t, movedFilePath, destFilePath)

				movedHostDescription, err := movedFile.GetHostDescription()
				require.NoError(t, err)
				destHostDescription, err := destFile.GetHostDescription()
				require.NoError(t, err)

				require.EqualValues(t, movedHostDescription, destHostDescription)

				exists, err = fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				exists, err = destFile.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)
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
				const verbose bool = true
				ctx := getCtx()

				srcFile := getFileToTest(tt.implementationName)
				err := srcFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)
				defer srcFile.Delete(ctx, &filesoptions.DeleteOptions{})

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(ctx, &filesoptions.DeleteOptions{})
				destFile.Delete(ctx, &filesoptions.DeleteOptions{})

				require.True(t, mustutils.Must(srcFile.Exists(ctx)))
				require.False(t, mustutils.Must(destFile.Exists(ctx)))

				err = srcFile.CopyToFile(destFile, verbose)
				require.NoError(t, err)

				require.True(t, mustutils.Must(srcFile.Exists(ctx)))
				require.True(t, mustutils.Must(destFile.Exists(ctx)))

				require.EqualValues(t, tt.content, srcFile.MustReadAsString())

				require.EqualValues(t, tt.content, destFile.MustReadAsString())
			},
		)
	}
}

func TestFile_String(t *testing.T) {
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

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(ctx, &filesoptions.DeleteOptions{})

				path, err := toTest.GetPath()
				require.NoError(t, err)
				stringOutput := toTest.String()
				sprintf := fmt.Sprintf("'%s'", toTest)

				require.True(t, pathsutils.IsAbsolutePath(path))
				require.EqualValues(t, path, stringOutput)
				require.EqualValues(t, "'"+path+"'", sprintf)
			},
		)
	}
}
