package hostnamecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewHostnameCmd() *cobra.Command {
	const short = "Hostname management commands"

	cmd := &cobra.Command{
		Use:   "hostname",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` hostname <subcommand>`,
	}

	cmd.AddCommand(
		NewGetHostnameCmd(),
		NewSetHostnameCmd(),
	)

	return cmd
}
