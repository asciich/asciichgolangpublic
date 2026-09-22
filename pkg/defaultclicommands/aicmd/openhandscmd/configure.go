package openhandscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/openhandsutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewConfigureCmd() *cobra.Command {
	const short = "Configure openhands for a generic OpenAI-compatible LLM."

	cmd := &cobra.Command{
		Use:   "configure",
		Short: short,
		Long: short + `

Reads the API token from '` + openhandsutils.OPENAI_API_KEY_ENV_VAR + `' env var.
The base URL can be specified via '--openai-base-url' or the '` + openhandsutils.OPENAI_BASE_URL_ENV_VAR + `' env var.

Usage:
  ` + os.Args[0] + ` ai openhands configure --url http://<openhands-host>:<port> --name <profile-name> --model <model-name> --openai-base-url http://<llm-host>:<port>/v1
  ` + os.Args[0] + ` ai openhands configure --url http://<openhands-host>:<port> --name <profile-name> --model <model-name>  # reads OPENAI_BASE_URL from env

Examples:
  # Configure with explicit URLs for both OpenHands and LLM provider
  OPENAI_API_KEY="your-api-key" ` + os.Args[0] + ` ai openhands configure --verbose --url=http://localhost:8000 --name=local-ollama --model="openai/qwen3:14b" --openai-base-url="http://llm-host:11434/v1"

  # Configure using environment variable for LLM base URL
  OPENAI_API_KEY="your-api-key" OPENAI_BASE_URL="http://llm-host:11434/v1" ` + os.Args[0] + ` ai openhands configure --verbose --url=http://localhost:8000 --name=local-ollama --model="openai/qwen3:14b"
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if url == "" {
				logging.LogFatal("Please specify --url for the OpenHands instance.")
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if name == "" {
				logging.LogFatal("Please specify --name for the LLM profile.")
			}

			model, err := cmd.Flags().GetString("model")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if model == "" {
				logging.LogFatal("Please specify --model for the LLM.")
			}

			baseUrl, err := cmd.Flags().GetString("openai-base-url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			mustutils.Must0(openhandsutils.ConfigureLlmProfile(ctx, url, &openhandsutils.ConfigureLlmProfileOptions{
				ProfileName: name,
				Model:       model,
				BaseUrl:     baseUrl,
				ReadFromEnv: true,
			}))

			logging.LogGoodByCtxf(ctx, "LLM profile '%s' configured on openhands '%s'.", name, url)
		},
	}

	cmd.Flags().String("url", "", "URL to openhands. E.g: http://localhost:8000")
	cmd.Flags().String("name", "", "Name for the LLM profile.")
	cmd.Flags().String("model", "", "Model name to use. E.g: gpt-4, claude-3, etc.")
	cmd.Flags().String("openai-base-url", "", "Base URL for the OpenAI-compatible API. Can also be read from "+openhandsutils.OPENAI_BASE_URL_ENV_VAR+" env var.")

	return cmd
}
