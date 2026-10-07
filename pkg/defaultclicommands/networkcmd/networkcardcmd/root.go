package networkcardcmd

import "github.com/spf13/cobra"

func NewNetworkCardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "networkcard",
		Short: "Commands related to network cards.",
	}

	cmd.AddCommand(
		NewDeactivateOffloading(),
		NewIsOffloadingDeactivatedCmd(),
	)

	return cmd
}
