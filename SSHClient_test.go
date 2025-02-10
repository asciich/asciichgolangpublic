package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
)

func TestSshClient_SshClientIsCommandExecutor(t *testing.T) {
	var sshClient commandexecutor.CommandExecutor = MustGetSshClientByHostName("abc")
	require.NotNil(t, sshClient)

	require.EqualValues(
		t,
		"abc",
		sshClient.MustGetHostDescription(),
	)
}
