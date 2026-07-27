package chromacmd

import "github.com/spf13/cobra"

func NewChromaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chroma",
		Short: "Chroma vector database related commands.",
	}

	cmd.AddCommand(
		NewCheckReachableCmd(),
		NewIndexDocumentsCmd(),
		NewRunMcpServerCmd(),
		NewQueryDocumentsCmd(),
	)

	return cmd
}
