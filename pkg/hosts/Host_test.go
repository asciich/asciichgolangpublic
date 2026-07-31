package hosts_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/hosts"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestHostGetHostName(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"hostname.asciich.ch", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host, err := hosts.GetHostByHostname(tt.hostname)
				require.NoError(t, err)

				hostName, err := host.GetHostName()
				require.NoError(t, err)
				require.EqualValues(t, tt.hostname, hostName)
			},
		)
	}
}

func TestHostGetHostDescripion(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)
	tests := []struct {
		hostname          string
		expectedReachable bool
	}{
		{"host.example.com", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				host, err := hosts.GetHostByHostname(tt.hostname)
				require.NoError(t, err)

				hostDescription, err := host.GetHostDescription()
				require.NoError(t, err)
				require.EqualValues(t, tt.hostname, hostDescription)
			},
		)
	}
}

func TestHost_GetDirectoryByPath(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		hostname       string
		dirPath        string
		expectedExists bool
	}{
		{"localhost", "/home/", true},
		{"localhost", "/home/does_not_exist", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				host, err := hosts.GetHostByHostname(tt.hostname)
				require.NoError(t, err)

				directory, err := host.GetDirectoryByPath(ctx, tt.dirPath)
				require.NoError(t, err)

				_, ok := directory.(*commandexecutorfileoo.Directory)
				require.True(t, ok)

				exists, err := directory.Exists(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedExists, exists)
			},
		)
	}
}

// Connections to the local host should use Bash and not SSH by default.
func TestHost_LocalHostUsesBashCommandExecutorByDefault(t *testing.T) {
	host, err := hosts.GetLocalHost()
	require.NoError(t, err)

	commandExecutorHost, ok := host.(*hosts.CommandExecutorHost)
	require.True(t, ok)

	commandExecutor, err := commandExecutorHost.GetCommandExecutor()
	require.NoError(t, err)

	_, ok = commandExecutor.(*commandexecutorbashoo.BashService)
	require.True(t, ok)
}

func Test_HostIsCommandExecutor(t *testing.T) {
	var host commandexecutorinterfaces.CommandExecutor
	var err error

	host, err = hosts.GetHostByHostname("example.com")
	require.NoError(t, err)
	require.NotNil(t, host)
}
