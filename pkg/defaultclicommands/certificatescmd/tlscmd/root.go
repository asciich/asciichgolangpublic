package tlscmd

import "github.com/spf13/cobra"

func NewTlsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tls",
		Short: "TLS certificate related commands",
	}

	cmd.AddCommand(
		NewGetFromWebserverCmd(),
	)

	return cmd
}
