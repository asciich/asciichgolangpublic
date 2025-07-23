package commandexecutorsshclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/commandexecutor"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/sshutils/commandexecutorsshclient"
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
