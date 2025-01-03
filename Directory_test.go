package asciichgolangpublic

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDirectoryToTest(implementationName string) (directory Directory) {
	const verbose = true

	if implementationName == "localDirectory" {
		directory = MustGetLocalDirectoryByPath(
			TemporaryDirectories().MustCreateEmptyTemporaryDirectoryAndGetPath(verbose),
		)
	} else if implementationName == "localCommandExecutorDirectory" {
		directory = MustGetLocalCommandExecutorDirectoryByPath(
			TemporaryDirectories().MustCreateEmptyTemporaryDirectoryAndGetPath(verbose),
		)
	} else {
		LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	}

	return directory
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				subDir := dir.MustCreateSubDirectory("subdir", verbose)

				assert.NotEqualValues(
					dir.MustGetPath(),
					subDir.MustGetPath(),
				)

				parentDir := subDir.MustGetParentDirectory()

				assert.EqualValues(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				dir.MustWriteStringToFileInDirectory("hello_world", verbose, "test.txt")

				assert.EqualValues(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				dir.MustWriteStringToFileInDirectory(tt.content, verbose, "test.txt")

				assert.EqualValues(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				dir := getDirectoryToTest(tt.implementationName)
				defer dir.Delete(verbose)

				dir.MustWriteStringToFileInDirectory("1234\nabc\n", verbose, "test.txt")

				assert.EqualValues(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testDirectory := getDirectoryToTest(tt.implementationName)
				defer testDirectory.Delete(verbose)

				testDirectory.MustCreateSubDirectory("test1", verbose)
				test2 := testDirectory.MustCreateSubDirectory("test2", verbose)

				test2.MustCreateSubDirectory("a", verbose)
				test2.MustCreateSubDirectory("b", verbose)
				test2.MustCreateSubDirectory("c", verbose)

				subDirectoryList := testDirectory.MustListSubDirectoryPaths(
					&ListDirectoryOptions{
						Recursive:           false,
						ReturnRelativePaths: true,
						Verbose:             verbose,
					},
				)

				assert.Len(subDirectoryList, 2)
				assert.EqualValues("test1", subDirectoryList[0])
				assert.EqualValues("test2", subDirectoryList[1])

				subDirectoryList = testDirectory.MustListSubDirectoryPaths(
					&ListDirectoryOptions{
						Recursive:           true,
						ReturnRelativePaths: true,
						Verbose:             verbose,
					},
				)

				assert.Len(subDirectoryList, 5)
				assert.EqualValues("test1", subDirectoryList[0])
				assert.EqualValues("test2", subDirectoryList[1])
				assert.EqualValues("test2/a", subDirectoryList[2])
				assert.EqualValues("test2/b", subDirectoryList[3])
				assert.EqualValues("test2/c", subDirectoryList[4])

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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose = true

				testDirectory := getDirectoryToTest(tt.implementationName)
				defer testDirectory.Delete(verbose)

				testDirectory.MustCreateSubDirectory("test1", verbose)
				test2 := testDirectory.MustCreateSubDirectory("test2", verbose)

				test2.MustCreateSubDirectory("a", verbose)
				test2.MustCreateSubDirectory("b", verbose)
				test2.MustCreateSubDirectory("c", verbose)

				subDirectoryList := testDirectory.MustListSubDirectories(
					&ListDirectoryOptions{
						Recursive: false,
					},
				)

				assert.Len(subDirectoryList, 2)
				assert.EqualValues("test1", subDirectoryList[0].MustGetBaseName())
				assert.EqualValues("test2", subDirectoryList[1].MustGetBaseName())
				assert.EqualValues(testDirectory.MustGetLocalPath(), subDirectoryList[0].MustGetDirName())
				assert.EqualValues(testDirectory.MustGetLocalPath(), subDirectoryList[1].MustGetDirName())

				subDirectoryList = testDirectory.MustListSubDirectories(
					&ListDirectoryOptions{
						Recursive: true,
					},
				)

				assert.Len(subDirectoryList, 5)
				assert.EqualValues(subDirectoryList[0].MustGetBaseName(), "test1")
				assert.EqualValues(subDirectoryList[1].MustGetBaseName(), "test2")
				assert.EqualValues(subDirectoryList[2].MustGetBaseName(), "a")
				assert.EqualValues(subDirectoryList[3].MustGetBaseName(), "b")
				assert.EqualValues(subDirectoryList[4].MustGetBaseName(), "c")
				assert.EqualValues(subDirectoryList[0].MustGetDirName(), testDirectory.MustGetLocalPath())
				assert.EqualValues(subDirectoryList[1].MustGetDirName(), testDirectory.MustGetLocalPath())
				assert.EqualValues(subDirectoryList[2].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
				assert.EqualValues(subDirectoryList[3].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
				assert.EqualValues(subDirectoryList[4].MustGetDirName(), filepath.Join(testDirectory.MustGetLocalPath(), "test2"))
			},
		)
	}
}
