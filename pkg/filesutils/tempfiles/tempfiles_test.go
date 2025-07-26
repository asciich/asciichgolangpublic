package tempfiles_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_CreateTemporaryFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tmpFilePath, err := tempfiles.CreateTemporaryFile(context.TODO())
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(tmpFilePath, "/tmp/"))
		require.True(t, strings.HasPrefix(filepath.Base(tmpFilePath), "tempfile"))
	})
}

func Test_CreateNamedTemporaryFile(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		tmpFilePath, err := tempfiles.CreateNamedTemporaryFile(context.TODO(), "")
		require.Error(t, err)
		require.Empty(t, tmpFilePath)
	})

	t.Run("happy path", func(t *testing.T) {
		tmpFilePath, err := tempfiles.CreateNamedTemporaryFile(context.TODO(), "abc")
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(tmpFilePath, "/tmp/"))
		require.True(t, strings.HasPrefix(filepath.Base(tmpFilePath), "abc"))
	})
}
