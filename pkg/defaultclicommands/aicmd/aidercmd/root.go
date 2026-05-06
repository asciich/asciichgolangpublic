package aidercmd

import "github.com/spf13/cobra"

func NewAiderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aider",
		Short: "aider related commands",
	}

	cmd.AddCommand(
		NewBuildContainerImageCmd(),
		NewShowRunCommandCmd(),
	)

	return cmd
}
