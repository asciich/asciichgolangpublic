package uuidcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewUuidCmd() *cobra.Command {
	const short = "UUID related commands."

	cmd := &cobra.Command{
		Use:   "uuid",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` uuid uuid`,
	}

	cmd.AddCommand(
		NewGenerateCmd(),
	)

	return cmd
}
