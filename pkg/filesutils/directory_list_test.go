package filesutils_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestDirectory_ListSubDirectories tests Directory.ListSubDirectories method
func TestDirectory_ListSubDirectories(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativedirectoryoo"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				// Create parent directory
				parentDir := getTemporaryDirectoryToTest(tt.implementationName)
				defer parentDir.Delete(ctx, &filesoptions.DeleteOptions{})

				// Create a file in parent (required by ListFiles implementation)
				_, err := parentDir.CreateFileInDirectory(ctx, ".keep", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				// Create subdirectories
				subDir1, err := parentDir.CreateSubDirectory(ctx, "subdir1", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				subDir2, err := parentDir.CreateSubDirectory(ctx, "subdir2", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				// List subdirectories (non-recursive)
				subDirs, err := parentDir.ListSubDirectories(ctx, &parameteroptions.ListDirectoryOptions{})
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(subDirs), 2)

				// Verify subdirectories are listed
				foundSubDir1 := false
				foundSubDir2 := false
				for _, subDir := range subDirs {
					path, err := subDir.GetPath()
					require.NoError(t, err)
					subDir1Path, err := subDir1.GetPath()
					require.NoError(t, err)
					subDir2Path, err := subDir2.GetPath()
					require.NoError(t, err)
					if path == subDir1Path {
						foundSubDir1 = true
					}
					if path == subDir2Path {
						foundSubDir2 = true
					}
				}
				require.True(t, foundSubDir1, "subdir1 not found in list")
				require.True(t, foundSubDir2, "subdir2 not found in list")
			},
		)
	}
}

// TestDirectory_CopyContentToDirectory tests Directory.CopyContentToDirectory method
func TestDirectory_CopyContentToDirectory(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativedirectoryoo"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				// Create source directory
				srcDir := getTemporaryDirectoryToTest(tt.implementationName)
				defer srcDir.Delete(ctx, &filesoptions.DeleteOptions{})

				srcDirPath, err := srcDir.GetPath()
				require.NoError(t, err)

				// Create files in source directory using native file creation
				file1Path := filepath.Join(srcDirPath, "file1.txt")
				file1 := getFileToTest(tt.implementationName, file1Path)
				err = file1.WriteString(ctx, "content1", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				file2Path := filepath.Join(srcDirPath, "file2.txt")
				file2 := getFileToTest(tt.implementationName, file2Path)
				err = file2.WriteString(ctx, "content2", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// Create destination directory
				destDir := getTemporaryDirectoryToTest(tt.implementationName)
				defer destDir.Delete(ctx, &filesoptions.DeleteOptions{})

				// Copy content
				err = srcDir.CopyContentToDirectory(ctx, destDir)
				require.NoError(t, err)

				// Verify files were copied
				destFile1, err := destDir.GetFileInDirectory("file1.txt")
				require.NoError(t, err)
				content1, err := destFile1.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "content1", content1)

				destFile2, err := destDir.GetFileInDirectory("file2.txt")
				require.NoError(t, err)
				content2, err := destFile2.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "content2", content2)
			},
		)
	}
}
