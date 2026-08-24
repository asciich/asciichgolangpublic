package openhandscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/aiproviders/swisscommyai"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewConfigureSwisscomMyAICmd() *cobra.Command {
	const short = "Configure openhands for myAI LLM of Swisscom."

	cmd := &cobra.Command{
		Use:   "configure-swisscom-myai",
		Short: short,
		Long: short + `

Needs the env var '` + swisscommyai.API_KEY_ENV_VAR_NAME + `' set.

Usage:
  ` + os.Args[0] + ` ai openhands configure-swisscom-myai` + os.Args[0] + ` ai openhands configure-swisscom-myai` + swisscommyai.API_KEY_ENV_VAR_NAME + `=<YOUR_API_KEY> ` + os.Args[0] + ` ai openhands configure-swisscom-myai`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if url == "" {
				logging.LogFatal("Please specify --url.")
			}

			mustutils.Must0(swisscommyai.AddLlmProfileToOpenhands(ctx, url))

			logging.LogGoodByCtxf(ctx, "Swisscom myAI configured on openhands '%s'.", url)
		},
	}

	cmd.Flags().String("url", "", "URL to openhands. E.g: http://localhost:8000")

	return cmd
}
