package tailscalecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewTailscaleCmd() *cobra.Command {
	const short = "Tailscale related commands"

	cmd := &cobra.Command{
		Use:   "tailscale",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn tailscale tailscale`,
	}

	cmd.AddCommand(
		NewDockerClientInstructionsCmd(),
		NewExampleWebserverCmd(),
		NewHttpRequestCmd(),
	)

	return cmd
}
