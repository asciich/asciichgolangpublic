package truststorecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewTrustStoreCmd() *cobra.Command {
	const short = "Truststore related commands"

	cmd := &cobra.Command{
		Use:   "truststore",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` certificates truststore truststore`,
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
