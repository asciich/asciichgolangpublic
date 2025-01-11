package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSshClient_SshClientIsCommandExecutor(t *testing.T) {
	var sshClient CommandExecutor = MustGetSshClientByHostName("abc")
	require.NotNil(t, sshClient)

	assert.EqualValues(
		t,
		"abc",
		sshClient.MustGetHostDescription(),
	)
}
