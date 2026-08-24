package filescmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewFilesCmd() *cobra.Command {
	const short = "File and directory related commands"

	cmd := &cobra.Command{
		Use:   "files",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` files files`,
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
