package tlscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewTlsCmd() *cobra.Command {
	const short = "TLS certificate related commands"

	cmd := &cobra.Command{
		Use:   "tls",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` certificates tls tls`,
	}

	cmd.AddCommand(
		NewGetFromWebserverCmd(),
		NewShowInfoCmd(),
	)

	return cmd
}
