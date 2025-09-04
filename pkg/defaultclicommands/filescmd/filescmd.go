package filescmd

import "github.com/spf13/cobra"

func NewFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "files",
		Short: "File and directory related commands",
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}