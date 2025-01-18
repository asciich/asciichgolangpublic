package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestDirectoryBase_SetAndGetParentDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				directoryBase := NewDirectoryBase()

				tempDir, err := os.MkdirTemp("", "test_direcotry")
				require.Nil(t, err)

				directory := MustGetLocalDirectoryByPath(tempDir)
				defer directory.Delete(verbose)

				directoryBase.MustSetParentDirectoryForBaseClass(directory)

				assert.EqualValues(
					directoryBase.MustGetParentDirectoryForBaseClass(),
					directory,
				)
			},
		)
	}
}

func TestDirectoryBase_ListFiles_withoutFilter(t *testing.T) {
	tests := []struct {
		fileImplementationToTest string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				directory := getDirectoryToTest(tt.fileImplementationToTest)

				defer directory.Delete(verbose)

				directory.MustCreateFileInDirectory(verbose, "a.txt")
				directory.MustCreateFileInDirectory(verbose, "a.log")
				directory.MustCreateFileInDirectory(verbose, "a.toc")
				directory.MustCreateFileInDirectory(verbose, "b.toc")

				fileList := directory.MustListFilePaths(
					&parameteroptions.ListFileOptions{
						ReturnRelativePaths: true,
					},
				)

				assert.EqualValues(
					[]string{"a.log", "a.toc", "a.txt", "b.toc"},
					fileList,
				)

			},
		)
	}
}

func TestDirectoryBase_ListFiles(t *testing.T) {
	tests := []struct {
		fileImplementationToTest string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				directory := getDirectoryToTest(tt.fileImplementationToTest)

				defer directory.Delete(verbose)

				directory.MustCreateFileInDirectory(verbose, "a.txt")
				directory.MustCreateFileInDirectory(verbose, "a.log")
				directory.MustCreateFileInDirectory(verbose, "a.toc")
				directory.MustCreateFileInDirectory(verbose, "b.toc")

				fileList := directory.MustListFilePaths(
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{".*.log", ".*.toc"},
						ReturnRelativePaths:  true,
					},
				)

				assert.EqualValues(
					[]string{"a.log", "a.toc", "b.toc"},
					fileList,
				)

			},
		)
	}
}

func TestDirectoryBase_DeleteFilesMatching(t *testing.T) {
	tests := []struct {
		fileImplementationToTest string
	}{
		{"localDirectory"},
		{"localCommandExecutorDirectory"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				directory := getDirectoryToTest(tt.fileImplementationToTest)

				txtFile := directory.MustCreateFileInDirectory(verbose, "a.txt")
				locFile := directory.MustCreateFileInDirectory(verbose, "a.log")
				tocFile := directory.MustCreateFileInDirectory(verbose, "a.toc")
				toc2File := directory.MustCreateFileInDirectory(verbose, "b.toc")

				directory.MustDeleteFilesMatching(
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{".*.log", ".*.toc"},
					},
				)

				assert.True(txtFile.MustExists(verbose))
				assert.False(locFile.MustExists(verbose))
				assert.False(tocFile.MustExists(verbose))
				assert.False(toc2File.MustExists(verbose))
			},
		)
	}
}
