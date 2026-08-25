package croncommandcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewCronCommandCmd() *cobra.Command {
	const short = "Run a command periodiacally as defined by a cron interval."

	cmd := &cobra.Command{
		Use:   "cron-command",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` shell croncommand cron-command`,
	}

	cmd.AddCommand(
		NewRunCmd(),
	)

	return cmd
}
