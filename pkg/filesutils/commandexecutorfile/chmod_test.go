package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func TestChmod(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Chmod(ctx, nil, "/tmp/test", &filesoptions.ChmodOptions{})
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Chmod(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.ChmodOptions{})
		require.Error(t, err)
	})
	t.Run("nil options returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Chmod(ctx, commandexecutorexecoo.Exec(), "/tmp/test", nil)
		require.Error(t, err)
	})
}

func TestGetAccessPermissions(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(nil, "/tmp")
		require.Error(t, err)
		require.Zero(t, perms)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Zero(t, perms)
	})
	t.Run("/tmp returns valid permissions", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(commandexecutorexecoo.Exec(), "/tmp")
		require.NoError(t, err)
		require.NotZero(t, perms)
	})
	t.Run("nonexistent file returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(commandexecutorexecoo.Exec(), "/nonexistent/file")
		require.Error(t, err)
		require.Zero(t, perms)
	})
}

func TestGetAccessPermissionsString(t *testing.T) {
	t.Run("/etc/hostname returns valid permissions string", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissionsString(commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.NotEmpty(t, perms)
	})
	t.Run("nonexistent file returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissionsString(commandexecutorexecoo.Exec(), "/nonexistent/file")
		require.Error(t, err)
		require.Empty(t, perms)
	})
}
