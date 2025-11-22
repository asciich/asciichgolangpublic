package atlassianconfluencecmd

import "github.com/spf13/cobra"

func NewAtlassianConfluenceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "atlassian-confluence",
		Short: "Commands for the Atlassian confluence wiki",
	}

	cmd.AddCommand(
		NewDownloadPageCmd(),
		NewGetChildPageIdsCmd(),
		NewGetPageContentCmd(),
		NewGetRequestCmd(),
	)

	return cmd
}
