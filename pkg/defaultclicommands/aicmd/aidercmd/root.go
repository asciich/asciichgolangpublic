package aidercmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewAiderCmd() *cobra.Command {
	const short = "aider related commands"

	cmd := &cobra.Command{
		Use:   "aider",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai aider aider`,
	}

	cmd.AddCommand(
		NewBuildContainerImageCmd(),
		NewShowRunCommandCmd(),
		NewRunAiderCmd(), // Add this line
	)

	return cmd
}
