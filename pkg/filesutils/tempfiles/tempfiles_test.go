package tempfiles_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_CreateTemporaryFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tmpFilePath, err := tempfiles.CreateTemporaryFile(context.TODO())
		defer func() { _ = nativefiles.Delete(context.TODO(), tmpFilePath, &filesoptions.DeleteOptions{}) }()
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(tmpFilePath, "/tmp/"))
		require.True(t, strings.HasPrefix(filepath.Base(tmpFilePath), "tempfile"))
	})
}

func Test_CreateNamedTemporaryFile(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		tmpFilePath, err := tempfiles.CreateNamedTemporaryFile(context.TODO(), "")
		defer func() { _ = nativefiles.Delete(context.TODO(), tmpFilePath, &filesoptions.DeleteOptions{}) }()
		require.Error(t, err)
		require.Empty(t, tmpFilePath)
	})

	t.Run("happy path", func(t *testing.T) {
		tmpFilePath, err := tempfiles.CreateNamedTemporaryFile(context.TODO(), "abc")
		defer func() { _ = nativefiles.Delete(context.TODO(), tmpFilePath, &filesoptions.DeleteOptions{}) }()
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(tmpFilePath, "/tmp/"))
		require.True(t, strings.HasPrefix(filepath.Base(tmpFilePath), "abc"))
	})
}

func TestCreateTemporaryFile(t *testing.T) {
	tests := []struct {
		content string
	}{
		{""},
		{"a"},
		{"hello world"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := context.TODO()

				tmpFilePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, tt.content)
				require.NoError(t, err)
				defer func() { _ = nativefiles.Delete(context.TODO(), tmpFilePath, &filesoptions.DeleteOptions{}) }()

				exists := nativefiles.Exists(ctx, tmpFilePath)
				require.True(t, exists)

				got, err := nativefiles.ReadAsString(ctx, tmpFilePath, &filesoptions.ReadOptions{})
				require.NoError(t, err)
				require.EqualValues(t, tt.content, got)
			},
		)
	}
}
