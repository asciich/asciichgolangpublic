package commandexecutorchecksumutils_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils/commandexecutorchecksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_GetMD5SumFromFileByPath(t *testing.T) {
	t.Run("empty file name", func(t *testing.T) {
		ctx := getCtx()
		got, err := commandexecutorchecksumutils.GetMD5SumFromFileByPath(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("nil commandExecutor", func(t *testing.T) {
		ctx := getCtx()
		got, err := commandexecutorchecksumutils.GetMD5SumFromFileByPath(ctx, nil, "/tmp/whatever")
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("empty file", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		got, err := commandexecutorchecksumutils.GetMD5SumFromFileByPath(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, "d41d8cd98f00b204e9800998ecf8427e", got)
	})

	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		got, err := commandexecutorchecksumutils.GetMD5SumFromFileByPath(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", got)
	})
}

func Test_GetSha1SumFromFile(t *testing.T) {
	t.Run("empty file name", func(t *testing.T) {
		ctx := getCtx()
		got, err := commandexecutorchecksumutils.GetSha1SumFromFile(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("empty file", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		got, err := commandexecutorchecksumutils.GetSha1SumFromFile(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", got)
	})

	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		got, err := commandexecutorchecksumutils.GetSha1SumFromFile(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", got)
	})
}

func Test_GetSha256SumFromFile(t *testing.T) {
	t.Run("empty file name", func(t *testing.T) {
		ctx := getCtx()
		got, err := commandexecutorchecksumutils.GetSha256SumFromFile(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("empty_file", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		const expected = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

		got, err := commandexecutorchecksumutils.GetSha256SumFromFile(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, expected, got)
	})

	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		const expected = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

		got, err := commandexecutorchecksumutils.GetSha256SumFromFile(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, expected, got)
	})

	t.Run("file size bigger than buffer", func(t *testing.T) {
		ctx := getCtx()

		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

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

		// The expected value was calculated using: seq $((3 * 1024)) | xargs -III echo hello world | sha256sum
		const expected = "df6930764673b9c9f62e319925cc5325d85d3284594f5fa6782a32c69e4da7cd"

		got, err := commandexecutorchecksumutils.GetSha256SumFromFile(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, expected, got)
	})
}

func Test_GetSha512SumFromFile(t *testing.T) {
	t.Run("empty file name", func(t *testing.T) {
		ctx := getCtx()
		got, err := commandexecutorchecksumutils.GetSha512SumFromFile(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("empty file", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		const expected = "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"

		got, err := commandexecutorchecksumutils.GetSha512SumFromFile(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, expected, got)
	})

	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()
		tempPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempPath, &filesoptions.DeleteOptions{})

		const expected = "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f"

		got, err := commandexecutorchecksumutils.GetSha512SumFromFile(ctx, commandexecutorexecoo.Exec(), tempPath)
		require.NoError(t, err)
		require.EqualValues(t, expected, got)
	})
}
