package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getDirectoryToTest(implementationName string) (directory Directory) {
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	if implementationName == "localDirectory" {
		return MustGetLocalDirectoryByPath(tempDir)
	}

	if implementationName == "localCommandExecutorDirectory" {
		return MustGetLocalCommandExecutorDirectoryByPath(tempDir)
	}

	logging.LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	err = os.Remove(tempDir)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nil
}

func TestDirectory_GetParentDirectory(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				subDir := dir.MustCreateSubDirectory("subdir", verbose)

				require.NotEqualValues(
					dir.MustGetPath(),
					subDir.MustGetPath(),
				)

				parentDir := subDir.MustGetParentDirectory()

				require.EqualValues(
					dir.MustGetPath(),
					parentDir.MustGetPath(),
				)
			},
		)
	}
}

func TestDirectory_ReadFileInDirectoryAsString(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				dir.MustWriteStringToFileInDirectory("hello_world", verbose, "test.txt")

				require.EqualValues(
					"hello_world",
					dir.MustReadFileInDirectoryAsString("test.txt"),
				)
			},
		)
	}
}

func TestDirectory_ReadFileInDirectoryAsInt64(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
		expectedInt64      int64
	}{
		{"localDirectory", "1234", 1234},
		{"localDirectory", "1234\n", 1234},
		{"localDirectory", "1234 ", 1234},
		{"localDirectory", " 1234", 1234},
		{"localDirectory", "\n1234\n", 1234},
		{"localDirectory", "\n1234", 1234},
		{"localCommandExecutorDirectory", "1234", 1234},
		{"localCommandExecutorDirectory", "1234\n", 1234},
		{"localCommandExecutorDirectory", "1234 ", 1234},
		{"localCommandExecutorDirectory", " 1234", 1234},
		{"localCommandExecutorDirectory", "\n1234\n", 1234},
		{"localCommandExecutorDirectory", "\n1234", 1234},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				dir.MustWriteStringToFileInDirectory(tt.content, verbose, "test.txt")

				require.EqualValues(
					tt.expectedInt64,
					dir.MustReadFileInDirectoryAsInt64("test.txt"),
				)
			},
		)
	}
}

func TestDirectory_ReadFirstLineOfFileInDirectoryAsString(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				dir.MustWriteStringToFileInDirectory("1234\nabc\n", verbose, "test.txt")

				require.EqualValues(
					"1234",
					dir.MustReadFirstLineOfFileInDirectoryAsString("test.txt"),
				)
			},
		)
	}
}

func TestDirectory_ListSubDirectories_RelativePaths(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				testDirectory := getDirectoryToTest(tt.implementationName)
				defer testDirectory.Delete(verbose)

				testDirectory.MustCreateSubDirectory("test1", verbose)
				test2 := testDirectory.MustCreateSubDirectory("test2", verbose)

				test2.MustCreateSubDirectory("a", verbose)
				test2.MustCreateSubDirectory("b", verbose)
				test2.MustCreateSubDirectory("c", verbose)

				subDirectoryList := testDirectory.MustListSubDirectoryPaths(
					&parameteroptions.ListDirectoryOptions{
						Recursive:           false,
						ReturnRelativePaths: true,
						Verbose:             verbose,
					},
				)

				require.Len(subDirectoryList, 2)
				require.EqualValues("test1", subDirectoryList[0])
				require.EqualValues("test2", subDirectoryList[1])

				subDirectoryList = testDirectory.MustListSubDirectoryPaths(
					&parameteroptions.ListDirectoryOptions{
						Recursive:           true,
						ReturnRelativePaths: true,
						Verbose:             verbose,
					},
				)

				require.Len(subDirectoryList, 5)
				require.EqualValues("test1", subDirectoryList[0])
				require.EqualValues("test2", subDirectoryList[1])
				require.EqualValues("test2/a", subDirectoryList[2])
				require.EqualValues("test2/b", subDirectoryList[3])
				require.EqualValues("test2/c", subDirectoryList[4])

			},
		)
	}

}

func TestDirectory_ListSubDirectories(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose = true

				testDirectory := getDirectoryToTest(tt.implementationName)
				defer testDirectory.Delete(verbose)

				testDirectory.MustCreateSubDirectory("test1", verbose)
				test2 := testDirectory.MustCreateSubDirectory("test2", verbose)

				test2.MustCreateSubDirectory("a", verbose)
				test2.MustCreateSubDirectory("b", verbose)
				test2.MustCreateSubDirectory("c", verbose)

				subDirectoryList := testDirectory.MustListSubDirectories(
					&parameteroptions.ListDirectoryOptions{
						Recursive: false,
					},
				)

				require.Len(subDirectoryList, 2)
				require.EqualValues("test1", subDirectoryList[0].MustGetBaseName())
				require.EqualValues("test2", subDirectoryList[1].MustGetBaseName())
				require.EqualValues(testDirectory.MustGetLocalPath(), subDirectoryList[0].MustGetDirName())
				require.EqualValues(testDirectory.MustGetLocalPath(), subDirectoryList[1].MustGetDirName())

				subDirectoryList = testDirectory.MustListSubDirectories(
					&parameteroptions.ListDirectoryOptions{
						Recursive: true,
					},
				)

				require.Len(subDirectoryList, 5)
				require.EqualValues(subDirectoryList[0].MustGetBaseName(), "test1")
				require.EqualValues(subDirectoryList[1].MustGetBaseName(), "test2")
				require.EqualValues(subDirectoryList[2].MustGetBaseName(), "a")
				require.EqualValues(subDirectoryList[3].MustGetBaseName(), "b")
				require.EqualValues(subDirectoryList[4].MustGetBaseName(), "c")
				require.EqualValues(subDirectoryList[0].MustGetDirName(), testDirectory.MustGetLocalPath())
				require.EqualValues(subDirectoryList[1].MustGetDirName(), testDirectory.MustGetLocalPath())
				require.EqualValues(subDirectoryList[2].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
				require.EqualValues(subDirectoryList[3].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
				require.EqualValues(subDirectoryList[4].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
			},
		)
	}
}
