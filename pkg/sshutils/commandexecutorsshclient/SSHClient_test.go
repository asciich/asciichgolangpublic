package commandexecutorsshclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/commandexecutorsshclient"
)

func TestSshClient_SshClientIsCommandExecutor(t *testing.T) {
	var sshClient commandexecutor.CommandExecutor
	var err error
	sshClient, err = commandexecutorsshclient.GetSshClientByHostName("abc")
	require.NotNil(t, sshClient)

	description, err := sshClient.GetHostDescription()
	require.NoError(t, err)

	require.EqualValues(t, "abc", description)
}
