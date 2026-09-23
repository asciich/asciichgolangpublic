package openwebuicmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewOpenWebUICmd() *cobra.Command {
	const short = "openwebui related commands"

	cmd := &cobra.Command{
		Use:   "openwebui",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai openwebui <command>`,
	}

	cmd.AddCommand(
		NewRunAsDockerContainerCmd(),
	)

	return cmd
}
