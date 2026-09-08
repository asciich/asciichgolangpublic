package commandexecutorsshclient_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/commandexecutorsshclient"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestSshClient_SshClientIsCommandExecutor(t *testing.T) {
	var sshClient commandexecutorinterfaces.CommandExecutor
	var err error
	sshClient, err = commandexecutorsshclient.GetSshClientByHostName("abc")
	require.NoError(t, err)
	require.NotNil(t, sshClient)

	description, err := sshClient.GetHostDescription()
	require.NoError(t, err)

	require.EqualValues(t, "abc", description)
}

func TestSshClient_IsRunningOnLocalhost(t *testing.T) {
	// Test that SSH client always returns false for IsRunningOnLocalhost,
	// even when the host is "localhost"
	testCases := []struct {
		name     string
		hostName string
	}{
		{"localhost", "localhost"},
		{"127.0.0.1", "127.0.0.1"},
		{"remote host", "remote.example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sshClient, err := commandexecutorsshclient.GetSshClientByHostName(tc.hostName)
			require.NoError(t, err)
			require.NotNil(t, sshClient)

			isLocalhost, err := sshClient.IsRunningOnLocalhost()
			require.NoError(t, err)
			require.False(t, isLocalhost, "SSH client should never be considered as running on localhost, even when host is '%s'", tc.hostName)
		})
	}
}

func TestSshClient_ConnectionRefusedOnClosedPort(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	ctx := getCtx()

	sshClient, err := commandexecutorsshclient.GetSshClientByHostName("localhost")
	require.NoError(t, err)
	require.NotNil(t, sshClient)

	require.NoError(t, sshClient.SetSshUserName("user"))

	// Pick a port that is (almost certainly) not open so the TCP connect is refused.
	const closedPort = 2

	if isOpen, err := netutils.IsTcpPortOpen(ctx, "localhost", closedPort); err == nil && isOpen {
		t.Skipf("Port %d unexpectedly open on localhost, cannot test connection refused", closedPort)
	}

	require.NoError(t, sshClient.SetSshPort(closedPort))

	_, err = sshClient.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"echo", "hello"},
		},
	)

	require.Error(t, err)
	require.True(t,
		netutilserrors.IsConnectionRefusedError(err),
		"expected a connection-refused error, got: %v", err,
	)
}
