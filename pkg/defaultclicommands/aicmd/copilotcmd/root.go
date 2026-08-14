package copilotcmd

import "github.com/spf13/cobra"

func NewCopilotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "copilot",
		Short: "Microsoft copilot related commands",
	}

	cmd.AddCommand(
		NewRunInCopilotCliContainerCmd(),
	)

	return cmd
}
