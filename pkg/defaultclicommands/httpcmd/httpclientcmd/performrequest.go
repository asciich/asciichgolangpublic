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

func NewPerformRequestCmd(options *httpclientcmdoptions.HttpClientCmdOptions) *cobra.Command {
	const short = "Perform a request and print response body to stdout. Use --method to specify the used in the request."

	cmd := &cobra.Command{
		Use:   "perform-request",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` http httpclient perform-request`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one URL")
			}

			url := normalizeUrl(args[0])

			if url == "" {
				logging.LogFatal("Please specify exactly one URL. Given argument is empty string")
			}

			method, err := cmd.Flags().GetString("method")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			client := options.GetClient()

			response := mustutils.Must(client.SendRequest(
				ctx,
				&httpoptions.RequestOptions{
					Url:    url,
					Method: method,
				},
			))

			fmt.Print(mustutils.Must(response.GetBodyAsString()))
		},
	}

	cmd.Flags().String("method", "GET", "HTTP method to perform.")

	return cmd
}
