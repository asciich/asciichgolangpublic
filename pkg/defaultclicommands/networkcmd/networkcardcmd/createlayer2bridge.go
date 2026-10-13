package networkcardcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils"
)

func NewCreateLayer2Bridge() *cobra.Command {
	const short = "Create a layer 2 Linux bridge (software switch) and attach the given network card to it. Useful to give VMs direct access to the LAN."

	cmd := &cobra.Command{
		Use:   "create-layer2-bridge <bridge-name> <network-card-name>",
		Short: short,
		Long: short + `

Creates the Linux bridge <bridge-name> and attaches <network-card-name> as bridge port.
Virtual machines (e.g. KVM/libvirt "bridged networking") attached to the bridge are directly connected to the
same LAN as the network card and get their IP addresses e.g. from the DHCP server of the LAN.

The bridge:
    - takes over the MAC address of <network-card-name>, so existing DHCP reservations keep working.
    - gets its IP configuration using DHCP. The IP of the host moves from the network card to the bridge.

The change is applied both:
    - running: systemd-networkd is reloaded and the network card is reconfigured.
      No reboot is needed. The command waits until the network card is attached to the bridge.
    - permanent: By writing the following systemd-networkd files:
        - '/etc/systemd/network/<bridge-name>.netdev'                              : Defines the bridge.
        - '/etc/systemd/network/05-<bridge-name>-bind-<network-card-name>.network' : Attaches the network card to the bridge.
        - '/etc/systemd/network/10-<bridge-name>.network'                          : IP configuration (DHCP) of the bridge.
      They are applied by systemd-networkd every boot.

The command is idempotent: An already up to date configuration and an already active bridge are left untouched.

Requires root privileges and only works on Linux with an active and enabled systemd-networkd.

Warning:
    - The network connection over <network-card-name> is shortly interrupted while it is moved into the bridge.
      Be careful when running this over SSH using the same network card.
    - Remove other systemd-networkd '.network' files matching <network-card-name> or <bridge-name> first.
      Otherwise they may take precedence or conflict with the generated configuration.

Usage:
    ` + os.Args[0] + ` network network-card create-layer2-bridge <bridge-name> <network-card-name>

Example:
    ` + os.Args[0] + ` network network-card create-layer2-bridge br0 enp0s25

Validate:
    ip -d link show br0      # Shows "bridge" as link type.
    bridge link show         # Lists enp0s25 with "master br0".
    networkctl status br0    # Shows the state and IP address of the bridge.
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 2 {
				logging.LogFatal("Please specify exactly a bridge name and a network card name: create-layer2-bridge <bridge-name> <network-card-name>")
			}

			bridgeName := args[0]
			networkCardName := args[1]

			mustutils.Must0(networkcardutils.CreateLayer2Bridge(ctx, bridgeName, networkCardName))

			logging.LogGoodByCtxf(ctx, "Layer 2 bridge '%s' with network card '%s' created.", bridgeName, networkCardName)
		},
	}

	return cmd
}
