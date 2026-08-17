package commandexecutorhost_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhost"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestCommandExecutorHost_HostnameOfLocalhost(t *testing.T) {
	host, err := hostsutils.GetLocalCommandExecutorHost()
	require.NoError(t, err)

	hostName, err := host.GetHostName()
	require.NoError(t, err)
	require.EqualValues(t, "localhost", hostName)
}

func Test_CommandExecutorHost_GetCpuArchitecture(t *testing.T) {
	ctx := getCtx()

	host, err := hostsutils.GetHostByHostname("localhost")
	require.NoError(t, err)

	commandExecutor := host.GetDeepCopyAsCommandExecutor()

	arch, err := commandExecutor.GetCPUArchitecture(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, arch)
}

func Test_CommandExecutorHost_GetCommandExecutor(t *testing.T) {
	t.Run("commandExecutor not set", func(t *testing.T) {
		commandExecutorHost := &commandexecutorhost.CommandExecutorHost{}

		commandExecutor, err := commandExecutorHost.GetCommandExecutor()
		require.Error(t, err)
		require.Nil(t, commandExecutor)
	})

	t.Run("commandExecutor set", func(t *testing.T) {
		commandExecutorHost := &commandexecutorhost.CommandExecutorHost{}
		err := commandExecutorHost.SetCommandExecutor(commandexecutorexecoo.Exec())
		require.NoError(t, err)

		commandExecutor, err := commandExecutorHost.GetCommandExecutor()
		require.NoError(t, err)
		require.NotNil(t, commandExecutor)
	})
}

func TestGetFileByPath(t *testing.T) {
	t.Run("no path given", func(t *testing.T) {
		commandExecutorHost := commandexecutorhost.NewCommandExecutorHost()
		err := commandExecutorHost.SetCommandExecutor(commandexecutorexecoo.Exec())
		require.NoError(t, err)

		got, err := commandExecutorHost.GetFileByPath("")
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("path given", func(t *testing.T) {
		commandExecutorHost := commandexecutorhost.NewCommandExecutorHost()
		err := commandExecutorHost.SetCommandExecutor(commandexecutorexecoo.Exec())
		require.NoError(t, err)

		got, err := commandExecutorHost.GetFileByPath("/home/file.txt")
		require.NoError(t, err)
		require.NotNil(t, got)

		path, err := got.GetPath()
		require.NoError(t, err)

		require.EqualValues(t, "/home/file.txt", path)
	})
}
