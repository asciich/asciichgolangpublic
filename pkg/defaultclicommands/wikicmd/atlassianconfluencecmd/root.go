package atlassianconfluencecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewAtlassianConfluenceCmd() *cobra.Command {
	const short = "Commands for the Atlassian confluence wiki"

	cmd := &cobra.Command{
		Use:   "atlassian-confluence",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` wiki atlassianconfluence atlassian-confluence`,
	}

	cmd.AddCommand(
		NewDownloadPageCmd(),
		NewGetChildPageIdsCmd(),
		NewGetPageContentCmd(),
		NewGetRequestCmd(),
	)

	return cmd
}
