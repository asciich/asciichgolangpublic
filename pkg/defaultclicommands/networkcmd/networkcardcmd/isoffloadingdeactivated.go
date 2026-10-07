package networkcardcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils"
)

func NewIsOffloadingDeactivatedCmd() *cobra.Command {
	const short = "Check if offloading is deactivated for given network card, both running and permanent."

	cmd := &cobra.Command{
		Use:   "is-offloading-deactivated <network-card-name>",
		Short: short,
		Long: short + `

Checks the following offloading features of the given network card:
    - TCP segmentation offload (TSO/TSO6)
    - Generic segmentation offload (GSO)
    - Generic receive offload (GRO)
    - Large receive offload (LRO)

Offloading is only considered deactivated if both:
    - running: None of the features above is currently enabled on the network card.
    - permanent: The systemd link file '/etc/systemd/network/10-<network-card-name>-no-offload.link' exists
      and sets all offload features to false.

This is the state ensured by the 'deactivate-offloading' command.

Prints 'yes' if offloading is deactivated, 'no' otherwise.
Any other output or a non zero exit code indicates an error during the check.

Does not require root privileges. Only works on Linux using systemd/udev.

Usage:
    ` + os.Args[0] + ` network network-card is-offloading-deactivated <network-card-name>

Example:
    ` + os.Args[0] + ` network network-card is-offloading-deactivated eno1
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one network card name.")
			}

			networkCardName := args[0]

			isDeactivated, err := networkcardutils.IsOffloadingDeactivated(ctx, networkCardName)
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if isDeactivated {
				logging.LogGoodByCtxf(ctx, "Offloading for network card '%s' is deactivated.", networkCardName)
				fmt.Println("yes")
			} else {
				logging.LogInfoByCtxf(ctx, "Offloading for network card '%s' is not (fully) deactivated.", networkCardName)
				fmt.Println("no")
			}
		},
	}

	return cmd
}
