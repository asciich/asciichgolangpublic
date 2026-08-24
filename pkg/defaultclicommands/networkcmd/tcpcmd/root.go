package tcpcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewTcpCmd() *cobra.Command {
	const short = "TCP related commands."

	cmd := &cobra.Command{
		Use:   "tcp",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network tcp tcp`,
	}

	cmd.AddCommand(
		NewIsPortOpenCmd(),
	)

	return cmd
}
