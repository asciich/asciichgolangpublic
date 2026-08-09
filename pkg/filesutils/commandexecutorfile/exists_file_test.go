package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func TestFileExists(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.FileExists(ctx, nil, "/etc/hostname")
		require.Error(t, err)
		require.False(t, exists)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.FileExists(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, exists)
	})
	t.Run("/etc/hostname exists", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.FileExists(ctx, commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.True(t, exists)
	})
	t.Run("nonexistent file returns false", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.FileExists(ctx, commandexecutorexecoo.Exec(), "/nonexistent/file")
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func TestDirectoryExists(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.DirectoryExists(ctx, nil, "/tmp")
		require.Error(t, err)
		require.False(t, exists)
	})
	t.Run("empty directoryPath returns error", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.DirectoryExists(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, exists)
	})
	t.Run("/tmp exists", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.DirectoryExists(ctx, commandexecutorexecoo.Exec(), "/tmp")
		require.NoError(t, err)
		require.True(t, exists)
	})
	t.Run("nonexistent directory returns false", func(t *testing.T) {
		ctx := getCtx()
		exists, err := commandexecutorfile.DirectoryExists(ctx, commandexecutorexecoo.Exec(), "/nonexistent/directory")
		require.NoError(t, err)
		require.False(t, exists)
	})
}
