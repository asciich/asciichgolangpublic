package truststorecmd

import "github.com/spf13/cobra"

func NewTrustStoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "truststore",
		Short: "Truststore related commands",
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
