package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func TestReadFirstNBytes(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		bytes, err := commandexecutorfile.ReadFirstNBytes(ctx, nil, "/etc/hostname", 10)
		require.Error(t, err)
		require.Nil(t, bytes)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()
		bytes, err := commandexecutorfile.ReadFirstNBytes(ctx, commandexecutorexecoo.Exec(), "", 10)
		require.Error(t, err)
		require.Nil(t, bytes)
	})
	t.Run("negative bytes returns error", func(t *testing.T) {
		ctx := getCtx()
		bytes, err := commandexecutorfile.ReadFirstNBytes(ctx, commandexecutorexecoo.Exec(), "/etc/hostname", -1)
		require.Error(t, err)
		require.Nil(t, bytes)
	})
	t.Run("/etc/hostname returns bytes", func(t *testing.T) {
		ctx := getCtx()
		bytes, err := commandexecutorfile.ReadFirstNBytes(ctx, commandexecutorexecoo.Exec(), "/etc/hostname", 10)
		require.NoError(t, err)
		require.Len(t, bytes, 10)
	})
}

func TestReadAsString(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		str, err := commandexecutorfile.ReadAsString(nil, "/etc/hostname")
		require.Error(t, err)
		require.Empty(t, str)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		str, err := commandexecutorfile.ReadAsString(commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Empty(t, str)
	})
	t.Run("/etc/hostname returns content", func(t *testing.T) {
		str, err := commandexecutorfile.ReadAsString(commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.NotEmpty(t, str)
	})
}

func TestReadAsBytes(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		bytes, err := commandexecutorfile.ReadAsBytes(nil, "/etc/hostname")
		require.Error(t, err)
		require.Nil(t, bytes)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		bytes, err := commandexecutorfile.ReadAsBytes(commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Nil(t, bytes)
	})
	t.Run("/etc/hostname returns content", func(t *testing.T) {
		bytes, err := commandexecutorfile.ReadAsBytes(commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.NotEmpty(t, bytes)
	})
}
