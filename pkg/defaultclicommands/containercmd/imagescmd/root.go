package imagescmd

import "github.com/spf13/cobra"

func NewImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "Handle docker images.",
	}

	cmd.AddCommand(
		NewDownloadAsArchiveCmd(),
	)

	return cmd
}
