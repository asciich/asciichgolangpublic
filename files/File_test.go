package files

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/testutils"
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
func getFileToTest(implementationName string) (file File) {
	if implementationName == "localFile" {
		return MustGetLocalFileByPath(createTempFileAndGetPath())
	}

	if implementationName == "localCommandExecutorFile" {
		return MustGetLocalCommandExecutorFileByPath(createTempFileAndGetPath())
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
				require := require.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				require.True(fileToTest.MustExists(verbose))

				fileToTest.MustDelete(verbose)

				require.False(fileToTest.MustExists(verbose))
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
				require := require.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				for i := 0; i < 10; i++ {
					fileToTest.MustTruncate(int64(i), verbose)
					require.EqualValues(
						fileToTest.MustGetSizeBytes(),
						int64(i),
					)
				}

				fileToTest.MustTruncate(0, verbose)

				require.EqualValues(
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
				require := require.New(t)

				const verbose bool = true

				fileToTest := getFileToTest(tt.implementationName)
				defer fileToTest.Delete(verbose)

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(verbose)
				destFile.Delete(verbose)

				require.True(fileToTest.MustExists(verbose))
				require.False(destFile.MustExists(verbose))

				movedFile := fileToTest.MustMoveToPath(destFile.MustGetPath(), false, verbose)

				require.EqualValues(
					movedFile.MustGetPath(),
					destFile.MustGetPath(),
				)

				require.EqualValues(
					movedFile.MustGetHostDescription(),
					destFile.MustGetHostDescription(),
				)

				require.False(fileToTest.MustExists(verbose))
				require.True(destFile.MustExists(verbose))
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
				require := require.New(t)

				const verbose bool = true

				srcFile := getFileToTest(tt.implementationName)
				srcFile.MustWriteString(tt.content, verbose)
				defer srcFile.Delete(verbose)

				destFile := getFileToTest(tt.implementationName)
				defer destFile.Delete(verbose)
				destFile.Delete(verbose)

				require.True(srcFile.MustExists(verbose))
				require.False(destFile.MustExists(verbose))

				srcFile.MustCopyToFile(destFile, verbose)

				require.True(srcFile.MustExists(verbose))
				require.True(destFile.MustExists(verbose))

				require.EqualValues(
					tt.content,
					srcFile.MustReadAsString(),
				)

				require.EqualValues(
					tt.content,
					destFile.MustReadAsString(),
				)
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

				toTest.MustChmod(
					&parameteroptions.ChmodOptions{
						PermissionsString: tt.permissionsString,
						Verbose:           verbose,
					},
				)

				require.EqualValues(
					t,
					tt.expectedPermissionString,
					toTest.MustGetAccessPermissionsString(),
				)

				require.EqualValues(
					t,
					unixfilepermissionsutils.MustGetPermissionsValue(tt.expectedPermissionString),
					toTest.MustGetAccessPermissions(),
				)
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

				path := toTest.MustGetPath()
				stringOutput := toTest.String()
				sprintf := fmt.Sprintf("'%s'", toTest)

				require.True(t, pathsutils.IsAbsolutePath(path))
				require.EqualValues(t, path, stringOutput)
				require.EqualValues(t, "'"+path+"'", sprintf)
			},
		)
	}
}
