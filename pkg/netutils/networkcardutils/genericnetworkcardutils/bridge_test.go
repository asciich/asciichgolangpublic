package genericnetworkcardutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils/genericnetworkcardutils"
)

func Test_CheckLayer2BridgeParameters(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		require.NoError(t, genericnetworkcardutils.CheckLayer2BridgeParameters("br0", "enp0s25"))
	})

	t.Run("empty bridge name", func(t *testing.T) {
		require.Error(t, genericnetworkcardutils.CheckLayer2BridgeParameters("", "enp0s25"))
	})

	t.Run("empty interface name", func(t *testing.T) {
		require.Error(t, genericnetworkcardutils.CheckLayer2BridgeParameters("br0", ""))
	})

	t.Run("same name", func(t *testing.T) {
		require.Error(t, genericnetworkcardutils.CheckLayer2BridgeParameters("br0", "br0"))
	})

	t.Run("name too long", func(t *testing.T) {
		require.Error(t, genericnetworkcardutils.CheckLayer2BridgeParameters("br0123456789abcd", "enp0s25"))
	})

	t.Run("shell unsafe name", func(t *testing.T) {
		require.Error(t, genericnetworkcardutils.CheckLayer2BridgeParameters("br0';reboot'", "enp0s25"))
	})
}

func Test_GetLayer2BridgeSystemdNetworkdFiles(t *testing.T) {
	t.Run("generates netdev, bind and bridge network", func(t *testing.T) {
		files, err := genericnetworkcardutils.GetLayer2BridgeSystemdNetworkdFiles("br0", "enp0s25", "D0:50:99:46:00:40")
		require.NoError(t, err)
		require.Len(t, files, 3)

		require.Equal(t, "/etc/systemd/network/br0.netdev", files[0].Path)
		require.Contains(t, files[0].Content, "[NetDev]\nName=br0\nKind=bridge\nMACAddress=d0:50:99:46:00:40\n")

		require.Equal(t, "/etc/systemd/network/05-br0-bind-enp0s25.network", files[1].Path)
		require.Contains(t, files[1].Content, "[Match]\nName=enp0s25\n\n[Network]\nBridge=br0\n")

		require.Equal(t, "/etc/systemd/network/10-br0.network", files[2].Path)
		require.Contains(t, files[2].Content, "[Match]\nName=br0\n\n[Network]\nDHCP=yes\n")
	})

	t.Run("invalid mac address", func(t *testing.T) {
		_, err := genericnetworkcardutils.GetLayer2BridgeSystemdNetworkdFiles("br0", "enp0s25", "not-a-mac")
		require.Error(t, err)
	})

	t.Run("zero mac address", func(t *testing.T) {
		_, err := genericnetworkcardutils.GetLayer2BridgeSystemdNetworkdFiles("br0", "enp0s25", "00:00:00:00:00:00")
		require.Error(t, err)
	})
}
