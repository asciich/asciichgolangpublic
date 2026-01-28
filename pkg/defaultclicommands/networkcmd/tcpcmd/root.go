package tcpcmd

import "github.com/spf13/cobra"

func NewTcpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tcp",
		Short: "TCP related commands.",
	}

	cmd.AddCommand(
		NewIsPortOpenCmd(),
	)

	return cmd
}