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

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Send get request and print response body to stdout.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one URL")
			}

			url := args[0]

			if url == "" {
				logging.LogFatal("Please specify exactly one URL. Given argument is empty string")
			}

			response := mustutils.Must(httputils.GetNativeClient().SendRequest(
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
