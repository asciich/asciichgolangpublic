package commandexecutorsshclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/commandexecutorsshclient"
)

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
