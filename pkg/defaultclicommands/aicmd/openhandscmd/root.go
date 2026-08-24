package openhandscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewOpenHandsCmd() *cobra.Command {
	const short = "openhands related commands"

	cmd := &cobra.Command{
		Use:   "openhands",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai openhands openhands`,
	}

	cmd.AddCommand(
		NewConfigureGoogleAIStudioCmd(),
		NewConfigureInfomaniakCmd(),
		NewConfigureSwisscomMyAICmd(),
		NewRunAsDockerContainerCmd(),
	)

	return cmd
}
