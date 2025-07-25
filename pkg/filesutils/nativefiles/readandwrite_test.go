package nativefiles_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_ReadAndWriteAsString(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		tmpPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tmpPath)

		err = nativefiles.WriteString(ctx, tmpPath, "hello world")
		require.NoError(t, err)

		content, err := nativefiles.ReadAsString(ctx, tmpPath)
		require.NoError(t, err)

		require.EqualValues(t, "hello world", content)
	})
}

func Test_ReadAndWriteAsBytes(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		tmpPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tmpPath)

		err = nativefiles.WriteBytes(ctx, tmpPath, []byte("hello world"))
		require.NoError(t, err)

		content, err := nativefiles.ReadAsBytes(ctx, tmpPath)
		require.NoError(t, err)

		require.EqualValues(t, []byte("hello world"), content)
	})
}
