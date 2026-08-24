package openhandscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/aiproviders/googleaistudio"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewConfigureGoogleAIStudioCmd() *cobra.Command {
	const short = "Configure openhands for Google AI Studio LLM."

	cmd := &cobra.Command{
		Use:   "configure-google-ai-studio",
		Short: short,
		Long: short + `

Needs the env var '` + googleaistudio.API_KEY_ENV_VAR_NAME + `' set.

Usage:
  ` + os.Args[0] + ` ai openhands configure-google-ai-studio` + os.Args[0] + ` ai openhands configure-google-ai-studio` + googleaistudio.API_KEY_ENV_VAR_NAME + `=<YOUR_API_KEY> ` + os.Args[0] + ` ai openhands configure-google-ai-studio`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if url == "" {
				logging.LogFatal("Please specify --url.")
			}

			mustutils.Must0(googleaistudio.AddLlmProfileToOpenhands(ctx, url))

			logging.LogGoodByCtxf(ctx, "Google AI Studio configured on openhands '%s'.", url)
		},
	}

	cmd.Flags().String("url", "", "URL to openhands. E.g: http://localhost:8000")

	return cmd
}
