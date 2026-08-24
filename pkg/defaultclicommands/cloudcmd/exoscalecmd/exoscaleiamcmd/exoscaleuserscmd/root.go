package exoscaleuserscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewUsersCmd() *cobra.Command {
	const short = "Manage exoscale users"

	cmd := &cobra.Command{
		Use:   "users",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` cloud exoscale exoscaleiam exoscaleusers users`,
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
