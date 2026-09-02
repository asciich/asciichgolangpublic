package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/fileinfo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestDirectory_GetFileInfoOfFilesInDirectory tests Directory.GetFileInfoOfFilesInDirectory method
func TestDirectory_GetFileInfoOfFilesInDirectory(t *testing.T) {
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

				// Create test files with known content
				file1, err := parentDir.CreateFileInDirectory(ctx, "file1.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file1.WriteString(ctx, "content1", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				file2, err := parentDir.CreateFileInDirectory(ctx, "file2.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file2.WriteString(ctx, "content2", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				file3, err := parentDir.CreateFileInDirectory(ctx, "file3.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file3.WriteString(ctx, "content3", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// Get file info for all files
				fileInfos, err := parentDir.GetFileInfoOfFilesInDirectory(ctx, &parameteroptions.ListFileOptions{})
				require.NoError(t, err)
				require.Len(t, fileInfos, 3)

				// Verify file info
				for _, fileInfo := range fileInfos {
					require.NotNil(t, fileInfo)

					path, err := fileInfo.GetPath()
					require.NoError(t, err)
					require.NotEmpty(t, path)

					sizeBytes, err := fileInfo.GetSizeBytes()
					require.NoError(t, err)
					require.Greater(t, sizeBytes, int64(0))
				}
			},
		)
	}
}

// TestDirectory_GetFileInfoOfFilesInDirectory_WithOptions tests Directory.GetFileInfoOfFilesInDirectory with various options
func TestDirectory_GetFileInfoOfFilesInDirectory_WithOptions(t *testing.T) {
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

				// Create test files with different patterns
				file1, err := parentDir.CreateFileInDirectory(ctx, "test1.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file1.WriteString(ctx, "content1", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				file2, err := parentDir.CreateFileInDirectory(ctx, "test2.log", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file2.WriteString(ctx, "content2", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				file3, err := parentDir.CreateFileInDirectory(ctx, "other.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file3.WriteString(ctx, "content3", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// Test with match pattern
				options := &parameteroptions.ListFileOptions{}
				err = options.SetMatchBasenamePattern([]string{"test.*"})
				require.NoError(t, err)

				fileInfos, err := parentDir.GetFileInfoOfFilesInDirectory(ctx, options)
				require.NoError(t, err)
				require.Len(t, fileInfos, 2)

				// Verify only matching files are returned
				for _, fileInfo := range fileInfos {
					path, err := fileInfo.GetPath()
					require.NoError(t, err)
					require.Contains(t, path, "test")
				}
			},
		)
	}
}

// TestDirectory_GetFileInfoOfFilesInDirectory_EmptyDirectory tests Directory.GetFileInfoOfFilesInDirectory on empty directory
func TestDirectory_GetFileInfoOfFilesInDirectory_EmptyDirectory(t *testing.T) {
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

				// Create empty parent directory
				parentDir := getTemporaryDirectoryToTest(tt.implementationName)
				defer parentDir.Delete(ctx, &filesoptions.DeleteOptions{})

				// Get file info for empty directory
				options := &parameteroptions.ListFileOptions{}
				options.SetAllowEmptyListIfNoFileIsFound(true)

				fileInfos, err := parentDir.GetFileInfoOfFilesInDirectory(ctx, options)
				require.NoError(t, err)
				require.Len(t, fileInfos, 0)
			},
		)
	}
}

// TestDirectory_GetFileInfoOfFilesInDirectory_NonRecursive tests Directory.GetFileInfoOfFilesInDirectory with non-recursive option
func TestDirectory_GetFileInfoOfFilesInDirectory_NonRecursive(t *testing.T) {
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

				// Create test files in parent directory
				file1, err := parentDir.CreateFileInDirectory(ctx, "file1.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file1.WriteString(ctx, "content1", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// Create subdirectory with files
				subDir, err := parentDir.CreateSubDirectory(ctx, "subdir", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				file2, err := subDir.CreateFileInDirectory(ctx, "file2.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				err = file2.WriteString(ctx, "content2", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// Test non-recursive listing
				options := &parameteroptions.ListFileOptions{}
				options.SetNonRecursive(true)

				fileInfos, err := parentDir.GetFileInfoOfFilesInDirectory(ctx, options)
				require.NoError(t, err)
				require.Len(t, fileInfos, 1) // Only file1.txt should be listed

				path, err := fileInfos[0].GetPath()
				require.NoError(t, err)
				require.Contains(t, path, "file1.txt")
			},
		)
	}
}

// TestDirectory_GetFileInfoOfFilesInDirectory_NilOptions tests Directory.GetFileInfoOfFilesInDirectory with nil options
func TestDirectory_GetFileInfoOfFilesInDirectory_NilOptions(t *testing.T) {
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

				// Test with nil options - should return error
				_, err := parentDir.GetFileInfoOfFilesInDirectory(ctx, nil)
				require.Error(t, err)
			},
		)
	}
}

// TestFileInfo_GetPath tests FileInfo.GetPath method
func TestFileInfo_GetPath(t *testing.T) {
	fileInfo := fileinfo.NewFileInfo()

	// Test with empty path
	_, err := fileInfo.GetPath()
	require.Error(t, err)

	// Test with valid path
	err = fileInfo.SetPath("/test/path.txt")
	require.NoError(t, err)

	path, err := fileInfo.GetPath()
	require.NoError(t, err)
	require.Equal(t, "/test/path.txt", path)
}

// TestFileInfo_GetSizeBytes tests FileInfo.GetSizeBytes method
func TestFileInfo_GetSizeBytes(t *testing.T) {
	fileInfo := fileinfo.NewFileInfo()

	// Test with default size (0)
	size, err := fileInfo.GetSizeBytes()
	require.NoError(t, err)
	require.Equal(t, int64(0), size)

	// Test with valid size
	err = fileInfo.SetSizeBytes(1024)
	require.NoError(t, err)

	size, err = fileInfo.GetSizeBytes()
	require.NoError(t, err)
	require.Equal(t, int64(1024), size)
}

// TestFileInfo_GetPathAndSizeBytes tests FileInfo.GetPathAndSizeBytes method
func TestFileInfo_GetPathAndSizeBytes(t *testing.T) {
	fileInfo := fileinfo.NewFileInfo()

	// Test with empty path
	_, _, err := fileInfo.GetPathAndSizeBytes()
	require.Error(t, err)

	// Test with valid path and size
	err = fileInfo.SetPath("/test/path.txt")
	require.NoError(t, err)

	err = fileInfo.SetSizeBytes(2048)
	require.NoError(t, err)

	path, size, err := fileInfo.GetPathAndSizeBytes()
	require.NoError(t, err)
	require.Equal(t, "/test/path.txt", path)
	require.Equal(t, int64(2048), size)
}
