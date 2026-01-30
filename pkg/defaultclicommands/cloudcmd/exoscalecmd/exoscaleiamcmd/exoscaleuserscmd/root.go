package exoscaleuserscmd

import "github.com/spf13/cobra"

func NewUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "users",
		Short: "Manage exoscale users",
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}