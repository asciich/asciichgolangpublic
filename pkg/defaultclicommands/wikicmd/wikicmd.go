package wikicmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/wikicmd/atlassianconfluencecmd"
	"os"
)

func NewWikiCmd() *cobra.Command {
	const short = "wiki related commands"

	cmd := &cobra.Command{
		Use:   "wiki",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` wiki wiki`,
	}

	cmd.AddCommand(
		atlassianconfluencecmd.NewAtlassianConfluenceCmd(),
	)

	return cmd
}
