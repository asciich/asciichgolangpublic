package checksumutils_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestChecksumsGetSha256SumFromString(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha256SumFromString(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetSha256SumFromBytes(t *testing.T) {
	tests := []struct {
		input            []byte
		expectedChecksum string
	}{
		{[]byte(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{[]byte("hello world"), "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha256SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetSha1SumFromString(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"hello world", "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha1SumFromString(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetSha1SumFromBytes(t *testing.T) {
	tests := []struct {
		input            []byte
		expectedChecksum string
	}{
		{[]byte(""), "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{[]byte("hello world"), "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetSha1SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetMD5SumFromString(t *testing.T) {
	tests := []struct {
		input            string
		expectedChecksum string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"hello world", "5eb63bbbe01eeed093cb22bb8f5acdc3"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetMD5SumFromString(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func TestChecksumsGetMD5SumFromBytes(t *testing.T) {
	tests := []struct {
		input            []byte
		expectedChecksum string
	}{
		{[]byte(""), "d41d8cd98f00b204e9800998ecf8427e"},
		{[]byte("hello world"), "5eb63bbbe01eeed093cb22bb8f5acdc3"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				calculatedSum := checksumutils.GetMD5SumFromBytes(tt.input)
				require.EqualValues(t, tt.expectedChecksum, calculatedSum)
			},
		)
	}
}

func Test_GetMD5SumFromFileByPath(t *testing.T) {
	t.Run("empty file name", func(t *testing.T) {
		ctx := getCtx()
		got, err := checksumutils.GetMD5SumFromFileByPath(ctx, "")
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("Non existing file", func(t *testing.T) {
		ctx := getCtx()
		got, err := checksumutils.GetMD5SumFromFileByPath(ctx, "/this/file/does/not/exist")
		require.Error(t, err)
		require.Empty(t, got)
		require.True(t, filesgeneric.IsErrFileNotFound(err))
	})

	t.Run("empty file", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		got, err := checksumutils.GetMD5SumFromFileByPath(ctx, tempPath)
		require.NoError(t, err)
		require.EqualValues(t, got, "d41d8cd98f00b204e9800998ecf8427e")
	})

	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		got, err := checksumutils.GetMD5SumFromFileByPath(ctx, tempPath)
		require.NoError(t, err)
		require.EqualValues(t, got, "5eb63bbbe01eeed093cb22bb8f5acdc3")
	})
}

func Test_GetSha256SumFromFile(t *testing.T) {
	t.Run("empty_file", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		const expected = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

		got, err := checksumutils.GetSha256SumFromFile(ctx, tempPath)
		require.NoError(t, err)

		require.EqualValues(t, expected, got)
	})

	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		const expected = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

		got, err := checksumutils.GetSha256SumFromFile(ctx, tempPath)
		require.NoError(t, err)

		require.EqualValues(t, expected, got)
	})

	t.Run("file size bigger than buffer", func(t *testing.T) {
		ctx := getCtx()

		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		//defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		fd, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		require.NoError(t, err)
		defer fd.Close()

		for range 3 * 1024 {
			fmt.Fprint(fd, "hello world\n")
		}
		err = fd.Close()
		require.NoError(t, err)

		size, err := nativefiles.GetSizeBytes(ctx, tempPath)
		require.NoError(t, err)

		require.Greater(t, size, int64(32*1024))

		// The expected value was caluculated using: seq $((3 * 1024)) | xargs -III echo hello world | sha256sum
		const expected = "df6930764673b9c9f62e319925cc5325d85d3284594f5fa6782a32c69e4da7cd"

		got, err := checksumutils.GetSha256SumFromFile(ctx, tempPath)
		require.NoError(t, err)
		require.EqualValues(t, expected, got)
	})
}
