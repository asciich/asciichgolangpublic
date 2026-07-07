package tempfiles_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateTemporaryFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		tmpFilePath, err := tempfiles.CreateTemporaryFile(ctx)
		defer func() { _ = nativefiles.Delete(ctx, tmpFilePath, &filesoptions.DeleteOptions{}) }()
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(tmpFilePath, "/tmp/"))
		require.True(t, strings.HasPrefix(filepath.Base(tmpFilePath), "tempfile"))
	})
}

func Test_CreateNamedTemporaryFile(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		ctx := getCtx()
		tmpFilePath, err := tempfiles.CreateNamedTemporaryFile(ctx, "")
		defer func() { _ = nativefiles.Delete(ctx, tmpFilePath, &filesoptions.DeleteOptions{}) }()
		require.Error(t, err)
		require.Empty(t, tmpFilePath)
	})

	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()
		tmpFilePath, err := tempfiles.CreateNamedTemporaryFile(ctx, "abc")
		defer func() { _ = nativefiles.Delete(ctx, tmpFilePath, &filesoptions.DeleteOptions{}) }()
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(tmpFilePath, "/tmp/"))
		require.True(t, strings.HasPrefix(filepath.Base(tmpFilePath), "abc"))
	})
}

func TestCreateTemporaryFileFromContentString(t *testing.T) {
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
				ctx := getCtx()

				tmpFilePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, tt.content)
				require.NoError(t, err)
				defer func() { _ = nativefiles.Delete(ctx, tmpFilePath, &filesoptions.DeleteOptions{}) }()

				exists := nativefiles.Exists(ctx, tmpFilePath)
				require.True(t, exists)

				got, err := nativefiles.ReadAsString(ctx, tmpFilePath, &filesoptions.ReadOptions{})
				require.NoError(t, err)
				require.EqualValues(t, tt.content, got)
			},
		)
	}
}

func TestCreateTemporaryFileFromContentBytes(t *testing.T) {
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
				ctx := getCtx()

				tmpFilePath, err := tempfiles.CreateTemporaryFileFromContentBytes(ctx, []byte(tt.content))
				require.NoError(t, err)
				defer func() { _ = nativefiles.Delete(ctx, tmpFilePath, &filesoptions.DeleteOptions{}) }()

				exists := nativefiles.Exists(ctx, tmpFilePath)
				require.True(t, exists)

				got, err := nativefiles.ReadAsString(ctx, tmpFilePath, &filesoptions.ReadOptions{})
				require.NoError(t, err)
				require.EqualValues(t, tt.content, got)
			},
		)
	}
}
