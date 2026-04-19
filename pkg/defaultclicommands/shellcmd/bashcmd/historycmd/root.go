package historycmd

import "github.com/spf13/cobra"

func NewHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Bash history related commands.",
	}

	cmd.AddCommand(
		NewEnableEmmediateWriteCmd(),
		NewIncreaseSizeCmd(),
	)

	return cmd
}
