package exoscalednscmd

import "github.com/spf13/cobra"

func NewDnsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Exoscale DNS.",
	}

	cmd.AddCommand(
		NewCreateRecordWithPublicIp(),
		NewListDomainsCmd(),
	)

	return cmd
}
