package k9scmd

import "github.com/spf13/cobra"

func NewK9sCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k9s",
		Short: "k9s related commands.",
	}

	cmd.AddCommand(
		NewInstallK9sCmd(),
	)

	return cmd
}
