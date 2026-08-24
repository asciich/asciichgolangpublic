package commandexecutorfile_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_OpenAsWriteCloser(t *testing.T) {
	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer writeCloser.Close()

		_, err = fmt.Fprintf(writeCloser, "hello world")
		require.NoError(t, err)
		err = writeCloser.Close()
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})
}

func Test_WriteBytes(t *testing.T) {
	t.Run("nil commandExecutor", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, nil, "/tmp/test.txt", []byte("hello"), &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("empty path", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), "", []byte("hello"), &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("nil content", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), "/tmp/test.txt", nil, &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("write and read back hello world", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte("hello world")
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})

	t.Run("write binary content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsBytes(ctx, tempFile)
		require.NoError(t, err)
		require.EqualValues(t, content, got)
	})

	t.Run("write empty content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte{}
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsBytes(ctx, tempFile)
		require.NoError(t, err)
		require.EqualValues(t, content, got)
	})

	t.Run("write multiline content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte("line1\nline2\nline3\n")
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "line1\nline2\nline3\n", got)
	})
}
