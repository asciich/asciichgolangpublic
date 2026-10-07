package networkcardcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils"
)

func NewDeactivateOffloading() *cobra.Command {
	const short = "Deactivate offloading for given network card. Useful to avoid network card hangs."

	cmd := &cobra.Command{
		Use:   "deactivate-offloading <network-card-name>",
		Short: short,
		Long: short + `

Deactivates the following offloading features of the given network card:
    - TCP segmentation offload (TSO/TSO6)
    - Generic segmentation offload (GSO)
    - Generic receive offload (GRO)
    - Large receive offload (LRO)

The change is applied both:
    - running: Directly on the network card (same as 'ethtool -K <network-card-name> tso off gso off gro off lro off').
      No reboot, driver reload or link reset is needed.
    - permanent: By writing the systemd link file '/etc/systemd/network/10-<network-card-name>-no-offload.link'.
      It is applied by udev every time the network card shows up (e.g. after a reboot or driver reload).

This is a known workaround for network card hangs like the e1000e "Detected Hardware Unit Hang"
which results in a complete loss of network connectivity until the network card is reset.

The command is idempotent: Already deactivated features and an already up to date link file are left untouched.

Requires root privileges and only works on Linux using systemd/udev.

Usage:
    ` + os.Args[0] + ` network network-card deactivate-offloading <network-card-name>

Example:
    ` + os.Args[0] + ` network network-card deactivate-offloading eno1

Validate:
    ` + os.Args[0] + ` network network-card is-offloading-deactivated <network-card-name>
Alternative to validate
    ethtool -k eno1 | grep -E '^(tcp-segmentation-offload|generic-segmentation-offload|generic-receive-offload|large-receive-offload)'
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one network card name.")
			}

			networkCardName := args[0]

			mustutils.Must0(networkcardutils.DeactivateOffloading(ctx, networkCardName))

			logging.LogGoodByCtxf(ctx, "Deactivated offloading for network card '%s'.", networkCardName)
		},
	}

	return cmd
}
