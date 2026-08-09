package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func TestGetSha256Sum(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, nil, "/etc/hostname")
		require.Error(t, err)
		require.Empty(t, sum)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Empty(t, sum)
	})

	t.Run("/etc/hostname returns valid sha256", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.NotEmpty(t, sum)
		require.Len(t, sum, 64)
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, commandexecutorexecoo.Exec(), "/nonexistent/file")
		require.Error(t, err)
		require.Empty(t, sum)
	})
}
