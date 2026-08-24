package k9scmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewK9sCmd() *cobra.Command {
	const short = "k9s related commands."

	cmd := &cobra.Command{
		Use:   "k9s",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes k9s k9s`,
	}

	cmd.AddCommand(
		NewInstallK9sCmd(),
	)

	return cmd
}
