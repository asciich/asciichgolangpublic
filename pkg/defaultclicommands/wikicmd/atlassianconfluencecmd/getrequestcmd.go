package atlassianconfluencecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/wikiutils/atlassianconfluenceutils"
)

func NewGetRequestCmd() *cobra.Command {
	const shortDescription = "Perform a GET request with the header for atlassian confluence wiki (containing authentication) for the specified URL. Prints out the received body."

	cmd := &cobra.Command{
		Use:   "get-request",
		Short: shortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exatly one URL.")
			}

			url := args[0]

			fmt.Println(mustutils.Must(atlassianconfluenceutils.GetRequest(ctx, url)))

			logging.LogGoodByCtxf(ctx, "Get wiki page content of '%s' finished.", url)
		},
	}

	return cmd
}
