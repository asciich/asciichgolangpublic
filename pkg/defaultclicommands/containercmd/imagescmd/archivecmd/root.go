package archivecmd

import "github.com/spf13/cobra"

func NewArchiveCmd() *cobra.Command {
	const short = "Handle container image archives."

	cmd := &cobra.Command{
		Use:   "archive",
		Short: short,
		Long: short + `

An image archive can be downloaded using the "` + DOWNLOAD_AS_ARCHIVE_USE + `" command.

For docker users an already available image can be exported using:
  docker save -o my_image.tar image_name:tag
`,
	}

	cmd.AddCommand(
		NewAddFileCmd(),
		NewCreateSingleFileArchive(),
		NewDownloadAsArchiveCmd(),
		NewListFilesCmd(),
	)

	return cmd
}
