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

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), tempFile)
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
