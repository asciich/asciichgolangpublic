package hosts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandExecutorHost_HostnameOfLocalhost(t *testing.T) {
	host := MustGetLocalCommandExecutorHost()
	hostName, err := host.GetHostName()
	require.NoError(t, err)
	require.EqualValues(t, "localhost", hostName)
}
