package documentationcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewDocumentationCmd(rootCmd *cobra.Command) *cobra.Command {
	const short = "Commands for documentation."

	cmd := &cobra.Command{
		Use:   "documentation",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` documentation documentation`,
	}

	cmd.AddCommand(
		NewGenerateMarkdownCmd(rootCmd),
	)

	return cmd
}
