package xzutils_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/compressionutils/xzutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestDecompress_EmptyArchivePath(t *testing.T) {
	ctx := getCtx()

	err := xzutils.Decompress(ctx, "", "/tmp/output.img")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "archivePath")
}

func TestDecompress_EmptyOutputPath(t *testing.T) {
	ctx := getCtx()

	err := xzutils.Decompress(ctx, "/tmp/archive.xz", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "outputPath")
}

func TestDecompress_NonExistentArchive(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "output.img")

	err := xzutils.Decompress(ctx, "/tmp/nonexistent_archive_12345.xz", outputPath)
	assert.Error(t, err)
}

func TestDecompress_InvalidXzFile(t *testing.T) {
	ctx := getCtx()

	tmpDir := t.TempDir()
	invalidArchivePath := filepath.Join(tmpDir, "invalid.xz")
	outputPath := filepath.Join(tmpDir, "output.img")

	err := os.WriteFile(invalidArchivePath, []byte("this is not xz data"), 0644)
	require.NoError(t, err)

	err = xzutils.Decompress(ctx, invalidArchivePath, outputPath)
	assert.Error(t, err)
}

func TestDecompress_ValidXzFile(t *testing.T) {
	ctx := getCtx()

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "testfile.txt")
	compressedPath := filepath.Join(tmpDir, "testfile.txt.xz")
	outputPath := filepath.Join(tmpDir, "testfile_decompressed.txt")

	expectedContent := "Hello, this is a test file for xz decompression.\n"
	err := os.WriteFile(inputPath, []byte(expectedContent), 0644)
	require.NoError(t, err)

	err = xzutils.Compress(ctx, inputPath, compressedPath)
	require.NoError(t, err)

	err = xzutils.Decompress(ctx, compressedPath, outputPath)
	require.NoError(t, err)

	actualContent, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Equal(t, expectedContent, string(actualContent))
}

func TestDecompress_LargeFile(t *testing.T) {
	ctx := getCtx()

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "largefile.bin")
	compressedPath := filepath.Join(tmpDir, "largefile.bin.xz")
	outputPath := filepath.Join(tmpDir, "largefile_decompressed.bin")

	// Create a 1MB file with repeating pattern
	data := make([]byte, 1*1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	err := os.WriteFile(inputPath, data, 0644)
	require.NoError(t, err)

	err = xzutils.Compress(ctx, inputPath, compressedPath)
	require.NoError(t, err)

	err = xzutils.Decompress(ctx, compressedPath, outputPath)
	require.NoError(t, err)

	actualContent, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Equal(t, data, actualContent)
}

func TestCompress_EmptyInputPath(t *testing.T) {
	ctx := getCtx()

	err := xzutils.Compress(ctx, "", "/tmp/output.xz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inputPath")
}

func TestCompress_EmptyOutputPath(t *testing.T) {
	ctx := getCtx()

	err := xzutils.Compress(ctx, "/tmp/input.img", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "outputPath")
}

func TestCompress_NonExistentInputFile(t *testing.T) {
	ctx := getCtx()

	outputPath := filepath.Join(t.TempDir(), "output.xz")

	err := xzutils.Compress(ctx, "/tmp/nonexistent_input_12345.img", outputPath)
	assert.Error(t, err)
}

func TestCompress_ValidFile(t *testing.T) {
	ctx := getCtx()

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "testfile.txt")
	compressedPath := filepath.Join(tmpDir, "testfile.txt.xz")

	content := "Hello, this is a test file for xz compression.\n"
	err := os.WriteFile(inputPath, []byte(content), 0644)
	require.NoError(t, err)

	err = xzutils.Compress(ctx, inputPath, compressedPath)
	require.NoError(t, err)

	// Verify compressed file exists and is smaller or at least different
	compressedInfo, err := os.Stat(compressedPath)
	require.NoError(t, err)
	assert.Greater(t, compressedInfo.Size(), int64(0))

	// Verify we can decompress back to original content
	decompressedPath := filepath.Join(tmpDir, "testfile_decompressed.txt")
	err = xzutils.Decompress(ctx, compressedPath, decompressedPath)
	require.NoError(t, err)

	actualContent, err := os.ReadFile(decompressedPath)
	require.NoError(t, err)
	assert.Equal(t, content, string(actualContent))
}

func TestCompress_EmptyFile(t *testing.T) {
	ctx := getCtx()

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "empty.txt")
	compressedPath := filepath.Join(tmpDir, "empty.txt.xz")
	decompressedPath := filepath.Join(tmpDir, "empty_decompressed.txt")

	err := os.WriteFile(inputPath, []byte{}, 0644)
	require.NoError(t, err)

	err = xzutils.Compress(ctx, inputPath, compressedPath)
	require.NoError(t, err)

	compressedInfo, err := os.Stat(compressedPath)
	require.NoError(t, err)
	assert.Greater(t, compressedInfo.Size(), int64(0))

	err = xzutils.Decompress(ctx, compressedPath, decompressedPath)
	require.NoError(t, err)

	actualContent, err := os.ReadFile(decompressedPath)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, actualContent)
}

func TestCompressAndDecompress_RoundTrip(t *testing.T) {
	ctx := getCtx()

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "roundtrip.bin")
	compressedPath := filepath.Join(tmpDir, "roundtrip.bin.xz")
	outputPath := filepath.Join(tmpDir, "roundtrip_decompressed.bin")

	// Create file with random-like content
	data := make([]byte, 64*1024)
	for i := range data {
		data[i] = byte((i*7 + 13) % 256)
	}
	err := os.WriteFile(inputPath, data, 0644)
	require.NoError(t, err)

	err = xzutils.Compress(ctx, inputPath, compressedPath)
	require.NoError(t, err)

	err = xzutils.Decompress(ctx, compressedPath, outputPath)
	require.NoError(t, err)

	actualContent, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Equal(t, data, actualContent)
}
