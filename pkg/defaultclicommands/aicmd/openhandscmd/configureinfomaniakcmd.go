package openhandscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/aiproviders/infomaniak"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewConfigureInfomaniakCmd() *cobra.Command {
	const short = "Configure openhands for Infomaniak LLM."

	cmd := &cobra.Command{
		Use:   "configure-infomaniak",
		Short: short,
		Long: short + `

Needs the env var '` + infomaniak.API_KEY_ENV_VAR_NAME + `' set.

Usage:
  ` + os.Args[0] + ` ai openhands configure-infomaniak` + os.Args[0] + ` ai openhands configure-infomaniak` + infomaniak.API_KEY_ENV_VAR_NAME + `=<YOUR_API_KEY> ` + os.Args[0] + ` ai openhands configure-infomaniak`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if url == "" {
				logging.LogFatal("Please specify --url.")
			}

			productId, err := cmd.Flags().GetString("product-id")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if productId == "" {
				logging.LogFatal("Please specify --product-id.")
			}

			mustutils.Must0(infomaniak.AddLlmProfileToOpenhands(ctx, url, productId))

			logging.LogGoodByCtxf(ctx, "Infomaniak configured on openhands '%s'.", url)
		},
	}

	cmd.Flags().String("url", "", "URL to openhands. E.g: http://localhost:8000")
	cmd.Flags().String("product-id", "", "Infomaniak product ID. E.g: 12345")

	return cmd
}
