package exoscalednscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewDnsCmd() *cobra.Command {
	const short = "Exoscale DNS."

	cmd := &cobra.Command{
		Use:   "dns",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` cloud exoscale exoscaledns dns`,
	}

	cmd.AddCommand(
		NewCreateRecordWithPublicIp(),
		NewListDomainsCmd(),
	)

	return cmd
}
