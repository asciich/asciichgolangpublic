package networkcardutils_test

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils"
)

// There is no "Output:" comment on purpose:
// The example is compiled but not executed by "go test" since it changes the network configuration of the host.
func ExampleCreateLayer2Bridge() {
	ctx := context.TODO()

	// Create the bridge 'br0' and attach the physical network card 'enp0s25' to it:
	//   - Requires root privileges and an active and enabled systemd-networkd.
	//   - The bridge takes over the MAC address of 'enp0s25' and uses DHCP.
	//   - Applied instantly and persistent across reboots.
	//   - Idempotent: Safe to call multiple times.
	mustutils.Must0(networkcardutils.CreateLayer2Bridge(ctx, "br0", "enp0s25"))
}
