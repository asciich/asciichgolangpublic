package hosts_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/hosts"
)

func TestCommandExecutorHost_HostnameOfLocalhost(t *testing.T) {
	host, err := hosts.GetLocalCommandExecutorHost()
	require.NoError(t, err)

	hostName, err := host.GetHostName()
	require.NoError(t, err)
	require.EqualValues(t, "localhost", hostName)
}

func Test_CommandExecutorHost_GetCpuArchitecture(t *testing.T) {
	ctx := getCtx()

	host, err := hosts.GetHostByHostname("localhost")
	require.NoError(t, err)

	commandExecutor := host.GetDeepCopyAsCommandExecutor()

	arch, err := commandExecutor.GetCPUArchitecture(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, arch)
}
