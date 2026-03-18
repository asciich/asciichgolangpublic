package commandexecutorfile_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_OpenAsReadCloser(t *testing.T) {
	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		readCloser, err := commandexecutorfile.OpenAsReadCloser(ctx, commandexecutorbashoo.Bash(), tempFile)
		require.NoError(t, err)
		defer readCloser.Close()

		got, err := io.ReadAll(readCloser)
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})
}
