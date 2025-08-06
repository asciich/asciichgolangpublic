package files_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
				const verbose bool = true
				ctx := getCtx()

				directoryBase := files.NewDirectoryBase()

				tempDir, err := os.MkdirTemp("", "test_direcotry")
				require.Nil(t, err)

				directory := files.MustGetLocalDirectoryByPath(tempDir)
				defer directory.Delete(ctx, &filesoptions.DeleteOptions{})

				directoryBase.MustSetParentDirectoryForBaseClass(directory)

				require.EqualValues(
					t,
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
				ctx := getCtx()
				const verbose = true

				directory := getDirectoryToTest(tt.fileImplementationToTest)

				defer directory.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := directory.CreateFileInDirectory(ctx, "a.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = directory.CreateFileInDirectory(ctx, "a.log", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = directory.CreateFileInDirectory(ctx, "a.toc", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = directory.CreateFileInDirectory(ctx, "b.toc", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				fileList, err := directory.ListFilePaths(
					ctx,
					&parameteroptions.ListFileOptions{
						ReturnRelativePaths: true,
					},
				)
				require.NoError(t, err)

				require.EqualValues(
					t,
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
				ctx := getCtx()
				const verbose bool = true

				directory := getDirectoryToTest(tt.fileImplementationToTest)

				defer directory.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := directory.CreateFileInDirectory(ctx, "a.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = directory.CreateFileInDirectory(ctx, "a.log", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = directory.CreateFileInDirectory(ctx, "a.toc", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = directory.CreateFileInDirectory(ctx, "b.toc", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				fileList, err := directory.ListFilePaths(
					ctx,
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{".*.log", ".*.toc"},
						ReturnRelativePaths:  true,
					},
				)
				require.NoError(t, err)

				require.EqualValues(
					t,
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
				ctx := getCtx()

				const verbose bool = true

				directory := getDirectoryToTest(tt.fileImplementationToTest)

				txtFile, err := directory.CreateFileInDirectory(ctx, "a.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				locFile, err := directory.CreateFileInDirectory(ctx, "a.log", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				tocFile, err := directory.CreateFileInDirectory(ctx, "a.toc", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				toc2File, err := directory.CreateFileInDirectory(ctx, "b.toc", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				directory.DeleteFilesMatching(
					ctx,
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{".*.log", ".*.toc"},
					},
				)

				exists, err := txtFile.Exists(verbose)
				require.NoError(t, err)
				require.True(t, exists)

				exists, err = locFile.Exists(verbose)
				require.NoError(t, err)
				require.False(t, exists)

				exists, err = tocFile.Exists(verbose)
				require.NoError(t, err)
				require.False(t, exists)

				exists, err = toc2File.Exists(verbose)
				require.NoError(t, err)
				require.False(t, exists)
			},
		)
	}
}
