package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
)

func TestSshClient_SshClientIsCommandExecutor(t *testing.T) {
	var sshClient commandexecutor.CommandExecutor = MustGetSshClientByHostName("abc")
	require.NotNil(t, sshClient)

	description, err := sshClient.GetHostDescription()
	require.NoError(t, err)

	require.EqualValues(
		t,
		"abc",
		description,
	)
}
