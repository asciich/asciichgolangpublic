package openhandscmd

import (
	"github.com/spf13/cobra"
)

func NewOpenHandsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openhands",
		Short: "openhands related commands",
	}

	cmd.AddCommand(
		NewConfigureGoogleAIStudioCmd(),
		NewConfigureSwisscomMyAICmd(),
		NewRunAsDockerContainerCmd(),
	)

	return cmd
}
