package hosts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandExecutorHost_HostnameOfLocalhost(t *testing.T) {
	require := require.New(t)

	host := MustGetLocalCommandExecutorHost()

	require.EqualValues(
		"localhost",
		host.MustGetHostName(),
	)
}
