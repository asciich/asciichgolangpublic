package headscalecmd

import "github.com/spf13/cobra"

func NewHeadscaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "headscale",
		Short: "Headscale is an open source, self-hosted implementation of the Tailscale control server.",
	}

	cmd.AddCommand(
		NewMinimalConfigCmd(),
	)

	return cmd
}