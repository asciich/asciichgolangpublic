package copilotcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewCopilotCmd() *cobra.Command {
	const short = "Microsoft copilot related commands"

	cmd := &cobra.Command{
		Use:   "copilot",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai copilot copilot`,
	}

	cmd.AddCommand(
		NewRunInCopilotCliContainerCmd(),
	)

	return cmd
}
