package documentationcmd

import "github.com/spf13/cobra"

func NewDocumentationCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: "documentation",
		Short: "Commands for documentation.",
	}

	cmd.AddCommand(
		NewGenerateMarkdownCmd(rootCmd),
	)

	return cmd
}