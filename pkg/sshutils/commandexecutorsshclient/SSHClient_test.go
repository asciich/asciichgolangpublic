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
