package filesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_IsFile(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, ""))
	})

	t.Run("nonexisting file", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, "/this/file/does/not/exist"))
	})

	t.Run("existing file", func(t *testing.T) {
		require.True(t, filesutils.IsFile(ctx, "/etc/hosts"))
	})

	t.Run("directory", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, "/etc"))
	})

	t.Run("directory2", func(t *testing.T) {
		require.False(t, filesutils.IsFile(ctx, "/etc/"))
	})
}
