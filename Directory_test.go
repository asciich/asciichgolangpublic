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

				dir.MustWriteStringToFileInDirectory("1234", verbose, "test.txt")

				assert.EqualValues(
					1234,
					dir.MustReadFileInDirectoryAsInt64("test.txt"),
				)
			},
		)
	}
}
