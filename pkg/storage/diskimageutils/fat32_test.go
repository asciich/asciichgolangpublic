package diskimageutils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/storage/diskimageutils"
)

// --- CreateFat32Image tests ---

func TestCreateFat32Image_EmptyOutputPath(t *testing.T) {
	ctx := getCtx()

	err := diskimageutils.CreateFat32Image(ctx, "", &diskimageutils.CreateFat32Options{
		SizeBytes: 64 * 1024 * 1024,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "outputPath")
}

func TestCreateFat32Image_NilOptions(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "options")
}

func TestCreateFat32Image_TooSmall(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 1024 * 1024, // 1 MB, too small
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "33 MB")
}

func TestCreateFat32Image_Valid(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes:   64 * 1024 * 1024,
		VolumeLabel: "TESTIMG",
	})
	require.NoError(t, err)

	info, err := os.Stat(outputPath)
	require.NoError(t, err)
	assert.Equal(t, int64(64*1024*1024), info.Size())
}

func TestCreateFat32Image_DefaultVolumeLabel(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	info, err := os.Stat(outputPath)
	require.NoError(t, err)
	assert.Equal(t, int64(33*1024*1024), info.Size())
}

// --- Fat32ListFiles tests ---

func TestFat32ListFiles_EmptyImagePath(t *testing.T) {
	ctx := getCtx()

	result, err := diskimageutils.Fat32ListFiles(ctx, "")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "imagePath")
}

func TestFat32ListFiles_NonExistentImage(t *testing.T) {
	ctx := getCtx()

	result, err := diskimageutils.Fat32ListFiles(ctx, "/tmp/nonexistent_fat32_12345.img")
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestFat32ListFiles_EmptyImage(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "empty.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	files, err := diskimageutils.Fat32ListFiles(ctx, outputPath)
	require.NoError(t, err)
	assert.Equal(t, []string{"."}, files)
}

func TestFat32ListFiles_WithFiles(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "withfiles.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 64 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "hello.txt", []byte("hello world"))
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "readme.md", []byte("# README"))
	require.NoError(t, err)

	files, err := diskimageutils.Fat32ListFiles(ctx, outputPath)
	require.NoError(t, err)

	assert.Contains(t, files, ".")
	assert.Contains(t, files, "./hello.txt")
	assert.Contains(t, files, "./readme.md")
	assert.Len(t, files, 3)
}

func TestFat32ListFiles_WithSubdirectories(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "subdirs.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 64 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "config.txt", []byte("arm_64bit=1"))
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "overlays/readme.txt", []byte("overlay info"))
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "overlays/sub/deep.txt", []byte("deep file"))
	require.NoError(t, err)

	files, err := diskimageutils.Fat32ListFiles(ctx, outputPath)
	require.NoError(t, err)

	assert.Contains(t, files, ".")
	assert.Contains(t, files, "./config.txt")
	assert.Contains(t, files, "./overlays")
	assert.Contains(t, files, "./overlays/readme.txt")
	assert.Contains(t, files, "./overlays/sub")
	assert.Contains(t, files, "./overlays/sub/deep.txt")
	assert.Len(t, files, 6)
}

// --- Fat32WriteFile tests ---

func TestFat32WriteFile_EmptyImagePath(t *testing.T) {
	ctx := getCtx()

	err := diskimageutils.Fat32WriteFile(ctx, "", "test.txt", []byte("hello"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "imagePath")
}

func TestFat32WriteFile_EmptyFilePath(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "", []byte("hello"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "filePath")
}

func TestFat32WriteFile_NonExistentImage(t *testing.T) {
	ctx := getCtx()

	err := diskimageutils.Fat32WriteFile(ctx, "/tmp/nonexistent_fat32_12345.img", "test.txt", []byte("hello"))
	assert.Error(t, err)
}

func TestFat32WriteFile_SimpleFile(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "simple.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	content := []byte("Hello, FAT32 world!\n")
	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "hello.txt", content)
	require.NoError(t, err)

	readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "hello.txt")
	require.NoError(t, err)
	assert.Equal(t, content, readBack)
}

func TestFat32WriteFile_EmptyContent(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "empty_content.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "empty.txt", []byte{})
	require.NoError(t, err)

	readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "empty.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte{}, readBack)
}

func TestFat32WriteFile_LargeFile(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "large.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 64 * 1024 * 1024,
	})
	require.NoError(t, err)

	// Create a file larger than one cluster (4KB)
	content := make([]byte, 32*1024) // 32 KB
	for i := range content {
		content[i] = byte(i % 256)
	}

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "large.bin", content)
	require.NoError(t, err)

	readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "large.bin")
	require.NoError(t, err)
	assert.Equal(t, content, readBack)
}

func TestFat32WriteFile_InSubdirectory(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "subdir.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	content := []byte("file in subdirectory")
	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "mydir/test.txt", content)
	require.NoError(t, err)

	readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "mydir/test.txt")
	require.NoError(t, err)
	assert.Equal(t, content, readBack)
}

func TestFat32WriteFile_NestedSubdirectories(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "nested.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	content := []byte("deeply nested content")
	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "a/b/c/deep.txt", content)
	require.NoError(t, err)

	readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "a/b/c/deep.txt")
	require.NoError(t, err)
	assert.Equal(t, content, readBack)
}

func TestFat32WriteFile_MultipleFilesInSameDirectory(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "multi.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 64 * 1024 * 1024,
	})
	require.NoError(t, err)

	files := map[string][]byte{
		"file1.txt": []byte("content 1"),
		"file2.txt": []byte("content 2"),
		"file3.txt": []byte("content 3"),
	}

	for name, content := range files {
		err = diskimageutils.Fat32WriteFile(ctx, outputPath, name, content)
		require.NoError(t, err)
	}

	for name, expectedContent := range files {
		readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, name)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, readBack)
	}
}

func TestFat32WriteFile_MultipleFilesInSubdirectory(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "multisubdir.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 64 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "dir/a.txt", []byte("aaa"))
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "dir/b.txt", []byte("bbb"))
	require.NoError(t, err)

	readA, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "dir/a.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("aaa"), readA)

	readB, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "dir/b.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("bbb"), readB)
}

// --- Fat32ReadFile tests ---

func TestFat32ReadFile_EmptyImagePath(t *testing.T) {
	ctx := getCtx()

	result, err := diskimageutils.Fat32ReadFile(ctx, "", "test.txt")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "imagePath")
}

func TestFat32ReadFile_EmptyFilePath(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	result, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "filePath")
}

func TestFat32ReadFile_NonExistentFile(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	result, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "nonexistent.txt")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestFat32ReadFile_NonExistentDirectory(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	result, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "nodir/file.txt")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestFat32ReadFile_RoundTrip(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "roundtrip.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	expectedContent := []byte("This is a round-trip test for FAT32 read/write.\nLine 2.\n")
	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "roundtrip.txt", expectedContent)
	require.NoError(t, err)

	actualContent, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "roundtrip.txt")
	require.NoError(t, err)
	assert.Equal(t, expectedContent, actualContent)
}

// --- Fat32Exists tests ---

func TestFat32Exists_EmptyImagePath(t *testing.T) {
	ctx := getCtx()

	result, err := diskimageutils.Fat32Exists(ctx, "", "test.txt")
	assert.False(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "imagePath")
}

func TestFat32Exists_EmptyPath(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "test.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	result, err := diskimageutils.Fat32Exists(ctx, outputPath, "")
	assert.False(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "filePath")
}

func TestFat32Exists_FileExists(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "exists.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "hello.txt", []byte("hello"))
	require.NoError(t, err)

	exists, err := diskimageutils.Fat32Exists(ctx, outputPath, "hello.txt")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestFat32Exists_FileDoesNotExist(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "noexist.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	exists, err := diskimageutils.Fat32Exists(ctx, outputPath, "missing.txt")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestFat32Exists_DirectoryExists(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "direxists.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "mydir/file.txt", []byte("content"))
	require.NoError(t, err)

	exists, err := diskimageutils.Fat32Exists(ctx, outputPath, "mydir")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestFat32Exists_NestedFileExists(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "nestedexists.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "a/b/c.txt", []byte("nested"))
	require.NoError(t, err)

	exists, err := diskimageutils.Fat32Exists(ctx, outputPath, "a/b/c.txt")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = diskimageutils.Fat32Exists(ctx, outputPath, "a/b")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = diskimageutils.Fat32Exists(ctx, outputPath, "a")
	require.NoError(t, err)
	assert.True(t, exists)
}

// --- Integration / combined tests ---

func TestFat32_FullWorkflow(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "workflow.img")

	// Create image
	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes:   64 * 1024 * 1024,
		VolumeLabel: "BOOT",
	})
	require.NoError(t, err)

	// Write multiple files
	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "config.txt", []byte("arm_64bit=1\n"))
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "cmdline.txt", []byte("console=serial0,115200\n"))
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "overlays/i2c.dtbo", []byte{0x01, 0x02, 0x03, 0x04})
	require.NoError(t, err)

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "overlays/spi.dtbo", []byte{0x05, 0x06, 0x07, 0x08})
	require.NoError(t, err)

	// List all files
	files, err := diskimageutils.Fat32ListFiles(ctx, outputPath)
	require.NoError(t, err)

	assert.Contains(t, files, ".")
	assert.Contains(t, files, "./config.txt")
	assert.Contains(t, files, "./cmdline.txt")
	assert.Contains(t, files, "./overlays")
	assert.Contains(t, files, "./overlays/i2c.dtbo")
	assert.Contains(t, files, "./overlays/spi.dtbo")
	assert.Len(t, files, 6)

	// Read files back
	configContent, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "config.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("arm_64bit=1\n"), configContent)

	cmdlineContent, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "cmdline.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("console=serial0,115200\n"), cmdlineContent)

	i2cContent, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "overlays/i2c.dtbo")
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, i2cContent)

	spiContent, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "overlays/spi.dtbo")
	require.NoError(t, err)
	assert.Equal(t, []byte{0x05, 0x06, 0x07, 0x08}, spiContent)

	// Check existence
	exists, err := diskimageutils.Fat32Exists(ctx, outputPath, "config.txt")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = diskimageutils.Fat32Exists(ctx, outputPath, "nonexistent.txt")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestFat32_LargeFileMultipleClusters(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "largefile.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 64 * 1024 * 1024,
	})
	require.NoError(t, err)

	// Create a 128 KB file spanning multiple clusters (cluster size = 4KB)
	content := make([]byte, 128*1024)
	for i := range content {
		content[i] = byte((i*7 + 13) % 256)
	}

	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "bigfile.bin", content)
	require.NoError(t, err)

	readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "bigfile.bin")
	require.NoError(t, err)
	assert.Equal(t, content, readBack)
}

func TestFat32_CaseInsensitive(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "case.img")

	err := diskimageutils.CreateFat32Image(ctx, outputPath, &diskimageutils.CreateFat32Options{
		SizeBytes: 33 * 1024 * 1024,
	})
	require.NoError(t, err)

	content := []byte("case test content")
	err = diskimageutils.Fat32WriteFile(ctx, outputPath, "MyFile.TXT", content)
	require.NoError(t, err)

	// FAT32 short names are case-insensitive
	readBack, err := diskimageutils.Fat32ReadFile(ctx, outputPath, "MYFILE.TXT")
	require.NoError(t, err)
	assert.Equal(t, content, readBack)
}
