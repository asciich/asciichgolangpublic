package wikicmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/wikicmd/atlassianconfluencecmd"
)

func NewWikiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wiki",
		Short: "wiki related commands",
	}

	cmd.AddCommand(
		atlassianconfluencecmd.NewAtlassianConfluenceCmd(),
	)

	return cmd
}
