package httpclientcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/httpclientcmd/httpclientcmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewGetCmd(options *httpclientcmdoptions.HttpClientCmdOptions) *cobra.Command {
	const short = "Send get request and print response body to stdout."

	cmd := &cobra.Command{
		Use:   "get",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` http httpclient get`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one URL")
			}

			url := normalizeUrl(args[0])

			if url == "" {
				logging.LogFatal("Please specify exactly one URL. Given argument is empty string")
			}

			client := options.GetClient()
			response := mustutils.Must(client.SendRequest(
				ctx,
				&httpoptions.RequestOptions{
					Url:    url,
					Method: "GET",
				},
			))

			fmt.Print(mustutils.Must(response.GetBodyAsString()))
		},
	}

	return cmd
}
