package tailscalecmd

import "github.com/spf13/cobra"

func NewTailscaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tailscale",
		Short: "Tailscale related commands",
	}

	cmd.AddCommand(
		NewDockerClientInstructionsCmd(),
		NewExampleWebserverCmd(),
		NewHttpRequestCmd(),
	)

	return cmd
}
