package files_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
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
		return files.MustGetLocalFileByPath(createTempFileAndGetPath())
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				fileToTest.MustWriteString(tt.content, verbose)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				exists, err := fileToTest.Exists(verbose)
				require.NoError(t, err)
				require.True(t, exists)

				err = fileToTest.Delete(verbose)
				require.NoError(t, err)

				exists, err = fileToTest.Exists(verbose)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				fileToTest.MustWriteString(
					"this is a\nhello world\nexample text.\n",
					verbose,
				)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(verbose)
				destFile.Delete(verbose)

				exists, err := fileToTest.Exists(verbose)
				require.NoError(t, err)
				require.True(t, exists)

				exists, err = destFile.Exists(verbose)
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

				exists, err = fileToTest.Exists(verbose)
				require.NoError(t, err)
				require.False(t, exists)

				exists, err = destFile.Exists(verbose)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				srcFile := getFileToTest(tt.implementationName)
				srcFile.MustWriteString(tt.content, verbose)
				defer srcFile.Delete(verbose)

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(verbose)
				destFile.Delete(verbose)

				require.True(t, mustutils.Must(srcFile.Exists(verbose)))
				require.False(t, mustutils.Must(destFile.Exists(verbose)))

				err := srcFile.CopyToFile(destFile, verbose)
				require.NoError(t, err)

				require.True(t, mustutils.Must(srcFile.Exists(verbose)))
				require.True(t, mustutils.Must(destFile.Exists(verbose)))

				require.EqualValues(t, tt.content, srcFile.MustReadAsString())

				require.EqualValues(t, tt.content, destFile.MustReadAsString())
			},
		)
	}
}

func TestFile_Chmod(t *testing.T) {
	tests := []struct {
		implementationName       string
		permissionsString        string
		expectedPermissionString string
	}{
		{"localFile", "u=rw,g=,o=", "u=rw,g=,o="},
		{"localCommandExecutorFile", "u=rw,g=,o=", "u=rw,g=,o="},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

				err := toTest.Chmod(
					&parameteroptions.ChmodOptions{
						PermissionsString: tt.permissionsString,
						Verbose:           verbose,
					},
				)
				require.NoError(t, err)

				accessPermissionsString, err := toTest.GetAccessPermissionsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedPermissionString, accessPermissionsString)

				accessPermissions, err := toTest.GetAccessPermissions()
				require.NoError(t, err)
				require.EqualValues(t, unixfilepermissionsutils.MustGetPermissionsValue(tt.expectedPermissionString), accessPermissions)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

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
