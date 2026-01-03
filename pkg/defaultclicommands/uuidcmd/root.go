package uuidcmd

import "github.com/spf13/cobra"

func NewUuidCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid",
		Short: "UUID related commands.",
	}

	cmd.AddCommand(
		NewGenerateCmd(),
	)

	return cmd
}
