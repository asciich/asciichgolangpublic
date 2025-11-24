package latexcmd

import "github.com/spf13/cobra"

func NewLatexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latex",
		Short: "Latex related commands.",
		Long:  "Latex related commands.",
	}

	cmd.AddCommand(
		NewCleanupTemporaryFilesCmd(),
	)

	return cmd
}
