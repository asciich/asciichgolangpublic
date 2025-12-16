package clientcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewPerformRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perform-request",
		Short: "Perform a request and print response body to stdout. Use --method to specify the used in the request.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one URL")
			}

			url := args[0]

			if url == "" {
				logging.LogFatal("Please specify exactly one URL. Given argument is empty string")
			}

			method, err := cmd.Flags().GetString("method")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			response := mustutils.Must(httputils.GetNativeClient().SendRequest(
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
