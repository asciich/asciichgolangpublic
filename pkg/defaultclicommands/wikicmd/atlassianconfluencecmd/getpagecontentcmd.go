package atlassianconfluencecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/wikiutils/atlassianconfluenceutils"
)

func NewGetPageContentCmd() *cobra.Command {
	const shortDescription = "Get the page content of the specified wiki page. Use the full URL of the page to read."

	cmd := &cobra.Command{
		Use:   "get-page-content",
		Short: shortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			prettyPrint, err := cmd.Flags().GetBool("pretty-print")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if len(args) != 1 {
				logging.LogFatal("Please specify exatly one URL to get the wiki page content.")
			}

			url := args[0]

			fmt.Println(mustutils.Must(atlassianconfluenceutils.GetPageContent(
				ctx,
				url,
				&atlassianconfluenceutils.GetContentOptions{
					PrettyPrint: prettyPrint,
				},
			)))

			logging.LogGoodByCtxf(ctx, "Get wiki page content of '%s' finished.", url)
		},
	}

	cmd.PersistentFlags().Bool("pretty-print", false, "Pretty print received page content.")

	return cmd
}
