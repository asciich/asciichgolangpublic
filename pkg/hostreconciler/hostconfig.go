package hostreconciler

import (
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils/genericnetworkcardutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// HostConfigFileName is the name of the file inside a config directory
// which contains the host configuration.
const HostConfigFileName = "host.yaml"

// hostConfig represents the content of the host.yaml file.
type hostConfig struct {
	Hostname string        `yaml:"hostname"`
	Network  networkConfig `yaml:"network"`
}

// networkConfig represents the optional 'network' section of the host.yaml file.
type networkConfig struct {
	Interfaces []NetworkInterface `yaml:"interfaces"`
}

// NetworkInterface is the desired configuration of a network interface (network card).
//
// It represents one entry of 'network.interfaces' in the host.yaml file
// and is returned by Source.GetNetworkInterfaces.
type NetworkInterface struct {
	// Name of the network card, e.g. 'enp0s25'.
	Name string `yaml:"name"`

	// Optional name of the layer 2 bridge to create with this network card attached as bridge port, e.g. 'br0'.
	// No bridge is created if empty.
	Layer2Bridge string `yaml:"layer2_bridge"`

	// If true, TSO/GSO/GRO/LRO offloading is deactivated for this network card (running and permanent).
	// If false, the offloading settings are left untouched (they are NOT reactivated).
	DeactivateOffloading bool `yaml:"deactivate_offloading"`
}

const networkInterfaceConfigHint = "Make sure 'name', 'layer2_bridge' and 'deactivate_offloading' are set in the same list item:\n  interfaces:\n    - name: enp0\n      layer2_bridge: br0\n      deactivate_offloading: true"

// validate returns an error if the network interface config is incomplete or invalid.
func (n *NetworkInterface) validate() error {
	if n.Name == "" {
		return tracederrors.TracedErrorf(
			"Missing 'name' in network interface config (layer2_bridge: '%s', deactivate_offloading: '%t'). %s",
			n.Layer2Bridge, n.DeactivateOffloading, networkInterfaceConfigHint,
		)
	}

	err := genericnetworkcardutils.CheckInterfaceName(n.Name)
	if err != nil {
		return err
	}

	if n.Layer2Bridge != "" {
		err = genericnetworkcardutils.CheckLayer2BridgeParameters(n.Layer2Bridge, n.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// hasActions returns true if at least one change is configured for the network interface.
func (n *NetworkInterface) hasActions() bool {
	return n.Layer2Bridge != "" || n.DeactivateOffloading
}

// validateNetworkInterfaces returns an error if a network interface is invalid or the network interfaces conflict with each other.
// An empty list is valid.
func validateNetworkInterfaces(interfaces []NetworkInterface) error {
	interfaceNames := map[string]bool{}
	bridgeToInterface := map[string]string{}

	for i, iface := range interfaces {
		err := iface.validate()
		if err != nil {
			return tracederrors.TracedErrorf("Invalid entry '%d' in 'network.interfaces': %w", i, err)
		}

		if interfaceNames[iface.Name] {
			return tracederrors.TracedErrorf("Network interface '%s' is defined more than once in 'network.interfaces'.", iface.Name)
		}

		interfaceNames[iface.Name] = true

		if iface.Layer2Bridge == "" {
			continue
		}

		otherInterface, alreadyUsed := bridgeToInterface[iface.Layer2Bridge]
		if alreadyUsed {
			return tracederrors.TracedErrorf(
				"Layer 2 bridge '%s' is used by network interface '%s' and '%s'. Only one network interface per layer 2 bridge is supported.",
				iface.Layer2Bridge, otherInterface, iface.Name,
			)
		}

		bridgeToInterface[iface.Layer2Bridge] = iface.Name
	}

	// Nested bridges are not supported: A configured bridge must not be configured as network interface itself.
	for _, iface := range interfaces {
		if iface.Layer2Bridge != "" && interfaceNames[iface.Layer2Bridge] {
			return tracederrors.TracedErrorf(
				"Layer 2 bridge '%s' of network interface '%s' is itself configured as network interface. Nested bridges are not supported.",
				iface.Layer2Bridge, iface.Name,
			)
		}
	}

	return nil
}
