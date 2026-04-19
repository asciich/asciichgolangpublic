package croncommandcmd

import "github.com/spf13/cobra"

func NewCronCommandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cron-command",
		Short: "Run a command periodiacally as defined by a cron interval.",
	}

	cmd.AddCommand(
		NewRunCmd(),
	)

	return cmd
}