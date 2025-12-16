package clientcmd

import "github.com/spf13/cobra"

func NewClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "HTTP client functions",
	}

	cmd.AddCommand(
		NewGetCmd(),
		NewPerformRequestCmd(),
	)

	return cmd
}
