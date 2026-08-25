package chromacmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewChromaCmd() *cobra.Command {
	const short = "Chroma vector database related commands."

	cmd := &cobra.Command{
		Use:   "chroma",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai vectordatabase chroma chroma`,
	}

	cmd.AddCommand(
		NewCheckReachableCmd(),
		NewIndexDocumentsCmd(),
		NewRunMcpServerCmd(),
		NewQueryDocumentsCmd(),
	)

	return cmd
}
