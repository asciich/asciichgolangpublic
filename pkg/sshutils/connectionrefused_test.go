package sshutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/commandexecutorsshclient"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/nativesshclient"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// getSshClientByImplementationName returns an SSH client for the given
// implementation, configured to connect to the given host/port/user.
func getSshClientByImplementationName(t *testing.T, implementationName string, hostName string, port int, user string, password string) commandexecutorinterfaces.CommandExecutor {
	t.Helper()

	if implementationName == "commandExecutorSshClient" {
		client, err := commandexecutorsshclient.GetSshClientByHostName(hostName)
		require.NoError(t, err)
		require.NoError(t, client.SetSshUserName(user))
		require.NoError(t, client.SetSshPort(port))
		return client
	}

	if implementationName == "nativeSshClient" {
		client, err := nativesshclient.NewSshClientByHostName(hostName)
		require.NoError(t, err)
		require.NoError(t, client.SetSshUserName(user))
		require.NoError(t, client.SetSshPort(port))
		return client
	}

	t.Fatalf("Unknown implementation name '%s'", implementationName)
	return nil
}

func TestSshClient_ConnectionRefusedOnClosedPort(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorSshClient"},
		{"nativeSshClient"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				// Pick a port that is (almost certainly) not open so the TCP connect is refused.
				const closedPort = 2

				if isOpen, err := netutils.IsTcpPortOpen(ctx, "localhost", closedPort); err == nil && isOpen {
					t.Skipf("Port %d unexpectedly open on localhost, cannot test connection refused", closedPort)
				}

				sshClient := getSshClientByImplementationName(t, tt.implementationName, "localhost", closedPort, "user", "pass")
				require.NotNil(t, sshClient)

				_, err := sshClient.RunCommandAndGetStdoutAsString(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: []string{"echo", "hello"},
					},
				)
				require.Error(t, err)
				require.True(t,
					netutilserrors.IsConnectionRefusedError(err),
					"[%s] expected a connection-refused error, got: %v", tt.implementationName, err,
				)
			},
		)
	}
}
