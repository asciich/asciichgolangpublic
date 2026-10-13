package networkcardutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils/nativenetworkcardutils"
)

// CreateLayer2Bridge creates the layer 2 Linux bridge (software switch) 'bridgeName' on localhost
// and attaches the network interface 'interfaceName' as bridge port, e.g. to give VMs direct access to the LAN.
//
// The configuration is done using systemd-networkd, applied instantly and persistent across reboots.
// The bridge takes over the MAC address of 'interfaceName' and gets its IP configuration using DHCP.
//
// Requires root privileges and an active and enabled systemd-networkd.
// This function is idempotent.
func CreateLayer2Bridge(ctx context.Context, bridgeName string, interfaceName string) error {
	return nativenetworkcardutils.CreateLayer2Bridge(ctx, bridgeName, interfaceName)
}
