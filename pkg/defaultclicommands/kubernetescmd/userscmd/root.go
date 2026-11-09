package userscmd

import "github.com/spf13/cobra"

func NewUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Handle users in kubernetes.",
	}

	cmd.AddCommand(
		NewWhoAmICmd(),
	)

	return cmd
}
