package nativefiles_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_ReadAndWriteAsString(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		tmpPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tmpPath, &filesoptions.DeleteOptions{})

		err = nativefiles.WriteString(ctx, tmpPath, "hello world")
		require.NoError(t, err)

		content, err := nativefiles.ReadAsString(ctx, tmpPath, &filesoptions.ReadOptions{})
		require.NoError(t, err)

		require.EqualValues(t, "hello world", content)
	})
}

func Test_ReadAndWriteAsBytes(t *testing.T) {
	permRW := os.FileMode(0600)
	permRWX := os.FileMode(0700)
	permReadOnly := os.FileMode(0444)
	permDefault := os.FileMode(0644)

	tests := []struct {
		name          string
		content       []byte
		perm          *os.FileMode
		expectedError bool
	}{
		{
			name:          "happy path",
			content:       []byte("hello world"),
			perm:          nil,
			expectedError: false,
		},
		{
			name:          "with perm 0600 (rw-------)",
			content:       []byte("hello world"),
			perm:          &permRW,
			expectedError: false,
		},
		{
			name:          "with perm 0700 (rwx------)",
			content:       []byte("hello world"),
			perm:          &permRWX,
			expectedError: false,
		},
		{
			name:          "with perm 0444 (read-only) sets permission but does not fail on initial write",
			content:       []byte("hello world"),
			perm:          &permReadOnly,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := getCtx()
			tmpPath, err := tempfiles.CreateTemporaryFile(ctx)
			require.NoError(t, err)
			defer nativefiles.Delete(ctx, tmpPath, &filesoptions.DeleteOptions{})

			err = nativefiles.WriteBytes(ctx, tmpPath, tt.content, &filesoptions.WriteOptions{
				Perm: tt.perm,
			})
			if tt.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			content, err := nativefiles.ReadAsBytes(ctx, tmpPath)
			require.NoError(t, err)

			require.EqualValues(t, tt.content, content)

			if tt.perm == nil {
				info, err := os.Stat(tmpPath)
				require.NoError(t, err)
				require.Equal(t, permDefault, info.Mode().Perm())
			} else {
				info, err := os.Stat(tmpPath)
				require.NoError(t, err)
				require.Equal(t, *tt.perm, info.Mode().Perm())
			}
		})
	}
}

func Test_WriteBytes_CreatesParentDirectoriesRecursively(t *testing.T) {
	tests := []struct {
		name          string
		subPath       string
		content       []byte
		expectedError bool
	}{
		{
			name:          "single level parent directory",
			subPath:       "level1/file.txt",
			content:       []byte("hello world"),
			expectedError: false,
		},
		{
			name:          "two level nested parent directories",
			subPath:       "level1/level2/file.txt",
			content:       []byte("hello world"),
			expectedError: false,
		},
		{
			name:          "deeply nested parent directories",
			subPath:       "level1/level2/level3/level4/level5/file.txt",
			content:       []byte("hello world"),
			expectedError: false,
		},
		{
			name:          "no parent directory (flat file)",
			subPath:       "file.txt",
			content:       []byte("hello world"),
			expectedError: false,
		},
		{
			name:          "empty content with nested directories",
			subPath:       "level1/level2/empty.txt",
			content:       []byte{},
			expectedError: false,
		},
		{
			name:          "binary content with nested directories",
			subPath:       "level1/level2/binary.bin",
			content:       []byte{0x00, 0xFF, 0x1A, 0x2B},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := getCtx()

			// Create a unique temp base directory per test to avoid conflicts
			baseDir, err := os.MkdirTemp("", "write-bytes-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(baseDir)

			fullPath := filepath.Join(baseDir, tt.subPath)

			err = nativefiles.WriteBytes(ctx, fullPath, tt.content, &filesoptions.WriteOptions{})
			if tt.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify the file exists
			_, err = os.Stat(fullPath)
			require.NoError(t, err, "expected file to exist at path: %s", fullPath)

			// Verify all parent directories were created
			parentDir := filepath.Dir(fullPath)
			_, err = os.Stat(parentDir)
			require.NoError(t, err, "expected parent directory to exist at path: %s", parentDir)

			// Verify the written content can be read back correctly
			readContent, err := nativefiles.ReadAsBytes(ctx, fullPath)
			require.NoError(t, err)
			require.EqualValues(t, tt.content, readContent)
		})
	}
}

func Test_OpenAsReadCloser(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []byte
	}{
		{
			name:    "empty file",
			content: "",
			want:    []byte{},
		},
		{
			name:    "hello world",
			content: "hello world",
			want:    []byte("hello world"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := getCtx()

			tmpPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, tt.content)
			require.NoError(t, err)
			defer nativefiles.Delete(ctx, tmpPath, &filesoptions.DeleteOptions{})

			readCloser, err := nativefiles.OpenAsReadCloser(ctx, tmpPath)
			require.NoError(t, err)
			defer readCloser.Close()

			got, err := io.ReadAll(readCloser)
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func Test_OpenAsWriteCloser(t *testing.T) {
	ctx := getCtx()

	tmpPath, err := tempfiles.CreateTemporaryFile(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, tmpPath, &filesoptions.DeleteOptions{})

	writeCloser, err := nativefiles.OpenAsWriteCloser(ctx, tmpPath, &filesoptions.WriteOptions{})
	require.NoError(t, err)
	defer writeCloser.Close()

	_, err = fmt.Fprint(writeCloser, "hello world")
	require.NoError(t, err)
	err = writeCloser.Close()
	require.NoError(t, err)

	got, err := nativefiles.ReadAsString(ctx, tmpPath, &filesoptions.ReadOptions{})
	require.NoError(t, err)

	require.EqualValues(t, "hello world", got)
}
