package atlassianconfluencecmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/wikiutils/atlassianconfluenceutils"
)

func NewGetChildPageIdsCmd() *cobra.Command {
	const short = "Prints out all child ids of the given wiki page. Use the URL or Id to specify the page to query."

	cmd := &cobra.Command{
		Use:   "get-child-page-ids",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` wiki atlassianconfluence get-child-page-ids`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty one page URL query.")
			}

			url := args[0]

			for _, id := range mustutils.Must(atlassianconfluenceutils.GetChildPageIds(ctx, url, &atlassianconfluenceutils.GetChildPageOptions{
				Recursive: true,
			})) {
				fmt.Println(id)
			}

			logging.LogGoodByCtxf(ctx, "Get child page ids of wiki page '%s' finished.", url)
		},
	}

	return cmd
}
