package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getDirectoryToTest(implementationName string) (directory filesinterfaces.Directory) {
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	if implementationName == "localDirectory" {
		return files.MustGetLocalDirectoryByPath(tempDir)
	}

	if implementationName == "localCommandExecutorDirectory" {
		return mustutils.Must(files.GetLocalCommandExecutorDirectoryByPath(tempDir))
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
				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				subDir, err := dir.CreateSubDirectory("subdir", verbose)
				require.NoError(t, err)

				dirPath, err := dir.GetPath()
				require.NoError(t, err)

				subDirPath, err := subDir.GetPath()
				require.NoError(t, err)

				require.NotEqualValues(t, dirPath, subDirPath)

				parentDir, err := subDir.GetParentDirectory()
				require.NoError(t, err)

				parentDirPath, err := parentDir.GetPath()
				require.NoError(t, err)

				require.EqualValues(t, dirPath, parentDirPath)
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
				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				_, err := dir.WriteStringToFileInDirectory("hello_world", verbose, "test.txt")
				require.NoError(t, err)

				content, err := dir.ReadFileInDirectoryAsString("test.txt")
				require.NoError(t, err)

				require.EqualValues(t, "hello_world", content)
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
				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				_, err := dir.WriteStringToFileInDirectory(tt.content, verbose, "test.txt")
				require.NoError(t, err)

				content, err := dir.ReadFileInDirectoryAsInt64("test.txt")
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedInt64, content)
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
				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				_, err := dir.WriteStringToFileInDirectory("1234\nabc\n", verbose, "test.txt")
				require.NoError(t, err)

				content, err := dir.ReadFirstLineOfFileInDirectoryAsString("test.txt")
				require.NoError(t, err)

				require.EqualValues(t, "1234", content)
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
				const verbose = true

				testDirectory := getDirectoryToTest(tt.implementationName)
				defer testDirectory.Delete(verbose)

				_, err := testDirectory.CreateSubDirectory("test1", verbose)
				require.NoError(t, err)

				test2, err := testDirectory.CreateSubDirectory("test2", verbose)
				require.NoError(t, err)

				_, err = test2.CreateSubDirectory("a", verbose)
				require.NoError(t, err)
				_, err = test2.CreateSubDirectory("b", verbose)
				require.NoError(t, err)
				_, err = test2.CreateSubDirectory("c", verbose)
				require.NoError(t, err)

				subDirectoryList, err := testDirectory.ListSubDirectoryPaths(
					&parameteroptions.ListDirectoryOptions{
						Recursive:           false,
						ReturnRelativePaths: true,
						Verbose:             verbose,
					},
				)
				require.NoError(t, err)

				require.Len(t, subDirectoryList, 2)
				require.EqualValues(t, "test1", subDirectoryList[0])
				require.EqualValues(t, "test2", subDirectoryList[1])

				subDirectoryList, err = testDirectory.ListSubDirectoryPaths(
					&parameteroptions.ListDirectoryOptions{
						Recursive:           true,
						ReturnRelativePaths: true,
						Verbose:             verbose,
					},
				)
				require.NoError(t, err)

				require.Len(t, subDirectoryList, 5)
				require.EqualValues(t, "test1", subDirectoryList[0])
				require.EqualValues(t, "test2", subDirectoryList[1])
				require.EqualValues(t, "test2/a", subDirectoryList[2])
				require.EqualValues(t, "test2/b", subDirectoryList[3])
				require.EqualValues(t, "test2/c", subDirectoryList[4])

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
				const verbose = true

				testDirectory := getDirectoryToTest(tt.implementationName)
				defer testDirectory.Delete(verbose)

				_, err := testDirectory.CreateSubDirectory("test1", verbose)
				require.NoError(t, err)
				test2, err := testDirectory.CreateSubDirectory("test2", verbose)
				require.NoError(t, err)

				_, err = test2.CreateSubDirectory("a", verbose)
				require.NoError(t, err)
				_, err = test2.CreateSubDirectory("b", verbose)
				require.NoError(t, err)
				_, err = test2.CreateSubDirectory("c", verbose)
				require.NoError(t, err)

				subDirectoryList, err := testDirectory.ListSubDirectories(
					&parameteroptions.ListDirectoryOptions{
						Recursive: false,
					},
				)
				require.NoError(t, err)

				require.Len(t, subDirectoryList, 2)

				baseName, err := subDirectoryList[0].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, "test1", baseName)

				baseName, err = subDirectoryList[1].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, "test2", baseName)

				testDirLocalPath, err := testDirectory.GetLocalPath()
				require.NoError(t, err)

				dirName, err := subDirectoryList[0].GetDirName()
				require.NoError(t, err)
				require.EqualValues(t, testDirLocalPath, dirName)

				dirName, err = subDirectoryList[1].GetDirName()
				require.NoError(t, err)
				require.EqualValues(t, testDirLocalPath, dirName)

				subDirectoryList, err = testDirectory.ListSubDirectories(
					&parameteroptions.ListDirectoryOptions{
						Recursive: true,
					},
				)
				require.NoError(t, err)

				require.Len(t, subDirectoryList, 5)

				baseName, err = subDirectoryList[0].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, baseName, "test1")

				baseName, err = subDirectoryList[1].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, baseName, "test2")

				baseName, err = subDirectoryList[2].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, baseName, "a")

				baseName, err = subDirectoryList[3].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, baseName, "b")

				baseName, err = subDirectoryList[4].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, baseName, "c")

				testDirPath, err := testDirectory.GetLocalPath()
				require.NoError(t, err)

				dirName, err = subDirectoryList[0].GetDirName()
				require.NoError(t, err)
				require.EqualValues(t, dirName, testDirPath)

				dirName, err = subDirectoryList[1].GetDirName()
				require.NoError(t, err)
				require.EqualValues(t, dirName, testDirPath)

				dirName, err = subDirectoryList[2].GetDirName()
				require.NoError(t, err)
				require.EqualValues(t, dirName, filepath.Join(testDirPath, "test2"))

				dirName, err = subDirectoryList[3].GetDirName()
				require.NoError(t, err)
				require.EqualValues(t, dirName, filepath.Join(testDirPath, "test2"))

				dirName, err = subDirectoryList[4].GetDirName()
				require.NoError(t, err)
				require.EqualValues(t, dirName, filepath.Join(testDirPath, "test2"))
			},
		)
	}
}
