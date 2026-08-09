package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func TestTruncate(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Truncate(ctx, nil, "/tmp/test", 0)
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Truncate(ctx, commandexecutorexecoo.Exec(), "", 0)
		require.Error(t, err)
	})
	t.Run("negative size returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Truncate(ctx, commandexecutorexecoo.Exec(), "/tmp/test", -1)
		require.Error(t, err)
	})
}

func TestGetSizeBytes(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		size, err := commandexecutorfile.GetSizeBytes(ctx, nil, "/etc/hostname")
		require.Error(t, err)
		require.Zero(t, size)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		size, err := commandexecutorfile.GetSizeBytes(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Zero(t, size)
	})
	t.Run("/etc/hostname returns size", func(t *testing.T) {
		ctx := getCtx()
		size, err := commandexecutorfile.GetSizeBytes(ctx, commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.Greater(t, size, int64(0))
	})
}
