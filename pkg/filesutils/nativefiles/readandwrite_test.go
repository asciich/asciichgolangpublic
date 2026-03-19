package nativefiles_test

import (
	"fmt"
	"io"
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
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		tmpPath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tmpPath, &filesoptions.DeleteOptions{})

		err = nativefiles.WriteBytes(ctx, tmpPath, []byte("hello world"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		content, err := nativefiles.ReadAsBytes(ctx, tmpPath)
		require.NoError(t, err)

		require.EqualValues(t, []byte("hello world"), content)
	})
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
