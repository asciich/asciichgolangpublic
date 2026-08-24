package historycmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewHistoryCmd() *cobra.Command {
	const short = "Bash history related commands."

	cmd := &cobra.Command{
		Use:   "history",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` shell bash history history`,
	}

	cmd.AddCommand(
		NewEnableEmmediateWriteCmd(),
		NewIncreaseSizeCmd(),
	)

	return cmd
}
