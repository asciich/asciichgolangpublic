package asciichgolangpublic

import (
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
