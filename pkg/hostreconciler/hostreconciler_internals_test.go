package hostreconciler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// fakeSource is a Source returning fixed values for testing.
type fakeSource struct {
	hostname          string
	networkInterfaces []NetworkInterface
}

func (f *fakeSource) GetHostName(ctx context.Context) (string, error) {
	return f.hostname, nil
}

func (f *fakeSource) GetNetworkInterfaces(ctx context.Context) ([]NetworkInterface, error) {
	return f.networkInterfaces, nil
}

func Test_hostConfig_ParseNetwork(t *testing.T) {
	t.Run("network interfaces are parsed", func(t *testing.T) {
		content := `---
hostname: host.example.com
network:
  interfaces:
    - name: enp0
      layer2_bridge: br0
      deactivate_offloading: true
    - name: enp1
      deactivate_offloading: true
    - name: enp2
      layer2_bridge: br2
`
		config := hostConfig{}
		err := yaml.Unmarshal([]byte(content), &config)
		require.NoError(t, err)

		require.Equal(t, "host.example.com", config.Hostname)
		require.Equal(
			t,
			[]NetworkInterface{
				{Name: "enp0", Layer2Bridge: "br0", DeactivateOffloading: true},
				{Name: "enp1", DeactivateOffloading: true},
				{Name: "enp2", Layer2Bridge: "br2"},
			},
			config.Network.Interfaces,
		)
	})

	t.Run("network section is optional", func(t *testing.T) {
		config := hostConfig{}
		err := yaml.Unmarshal([]byte("---\nhostname: host.example.com\n"), &config)
		require.NoError(t, err)

		require.Empty(t, config.Network.Interfaces)
		require.NoError(t, validateNetworkInterfaces(config.Network.Interfaces))
	})

	t.Run("deactivate_offloading defaults to false", func(t *testing.T) {
		content := `---
interfaces:
  - name: enp0
    layer2_bridge: br0
`
		config := networkConfig{}
		err := yaml.Unmarshal([]byte(content), &config)
		require.NoError(t, err)

		require.Len(t, config.Interfaces, 1)
		require.False(t, config.Interfaces[0].DeactivateOffloading)
	})
}

func Test_validateNetworkInterfaces(t *testing.T) {
	t.Run("valid network interfaces", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", Layer2Bridge: "br0", DeactivateOffloading: true},
			{Name: "enp1", DeactivateOffloading: true},
			{Name: "enp2", Layer2Bridge: "br2"},
		}
		require.NoError(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("nil is valid", func(t *testing.T) {
		require.NoError(t, validateNetworkInterfaces(nil))
	})

	t.Run("network interface without actions is valid", func(t *testing.T) {
		require.NoError(t, validateNetworkInterfaces([]NetworkInterface{{Name: "enp0"}}))
	})

	t.Run("settings in separate list items", func(t *testing.T) {
		content := `---
interfaces:
  - name: enp0
  - layer2_bridge: br0
`
		config := networkConfig{}
		err := yaml.Unmarshal([]byte(content), &config)
		require.NoError(t, err)

		err = validateNetworkInterfaces(config.Interfaces)
		require.Error(t, err)
		require.Contains(t, err.Error(), "same list item")
	})

	t.Run("missing name", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{DeactivateOffloading: true}}))
	})

	t.Run("network interface defined twice", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", DeactivateOffloading: true},
			{Name: "enp0", Layer2Bridge: "br0"},
		}
		require.Error(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("bridge used by two network interfaces", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", Layer2Bridge: "br0"},
			{Name: "enp1", Layer2Bridge: "br0"},
		}
		require.Error(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("nested bridges", func(t *testing.T) {
		interfaces := []NetworkInterface{
			{Name: "enp0", Layer2Bridge: "br0"},
			{Name: "br0", Layer2Bridge: "br1"},
		}
		require.Error(t, validateNetworkInterfaces(interfaces))
	})

	t.Run("same name for bridge and network interface", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{Name: "br0", Layer2Bridge: "br0"}}))
	})

	t.Run("shell unsafe network interface name", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{Name: "enp0';reboot'", DeactivateOffloading: true}}))
	})

	t.Run("shell unsafe bridge name", func(t *testing.T) {
		require.Error(t, validateNetworkInterfaces([]NetworkInterface{{Name: "enp0", Layer2Bridge: "br0';reboot'"}}))
	})
}

func Test_reconcileNetworkInterfaces(t *testing.T) {
	ctx := context.TODO()

	t.Run("no network interfaces is a no-op without root", func(t *testing.T) {
		require.NoError(t, reconcileNetworkInterfaces(ctx, nil))
	})

	t.Run("network interface without actions is a no-op without root", func(t *testing.T) {
		require.NoError(t, reconcileNetworkInterfaces(ctx, []NetworkInterface{{Name: "enp0"}}))
	})

	t.Run("invalid network interfaces are rejected before any change", func(t *testing.T) {
		require.Error(t, reconcileNetworkInterfaces(ctx, []NetworkInterface{{Layer2Bridge: "br0"}}))
	})
}

func Test_ReconcileOnce(t *testing.T) {
	ctx := context.TODO()

	t.Run("nil source", func(t *testing.T) {
		require.Error(t, ReconcileOnce(ctx, nil))
	})

	t.Run("empty hostname", func(t *testing.T) {
		require.Error(t, ReconcileOnce(ctx, &fakeSource{}))
	})

	t.Run("invalid network interfaces are rejected before the hostname is changed", func(t *testing.T) {
		source := &fakeSource{
			hostname:          "host.example.com",
			networkInterfaces: []NetworkInterface{{Layer2Bridge: "br0"}},
		}

		err := ReconcileOnce(ctx, source)
		require.Error(t, err)

		// The validation error proves ReconcileOnce stopped before calling SetHostName.
		require.Contains(t, err.Error(), "same list item")
	})
}
