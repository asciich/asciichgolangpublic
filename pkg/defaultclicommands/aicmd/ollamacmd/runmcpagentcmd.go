package ollamacmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils/ollamaagent"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunMcpAgentCmd() *cobra.Command {
	const short = "Run an interactive ollama MCP agent with tool-calling support."

	cmd := &cobra.Command{
		Use:   "run-mcp-agent",
		Short: short,
		Long: short + `

The agent connects to one or more MCP servers, discovers their tools, and lets 
the ollama model decide when to call them.

Configuration can be provided either via a YAML config file or via CLI flags.
CLI flags take precedence over the config file.

Example using CLI flags:
  ` + os.Args[0] + ` ai ollama run-mcp-agent` + os.Args[0] + ` ai ollama run-mcp-agent`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			ollamaUrl, err := cmd.Flags().GetString("ollama-url")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			llmModel, err := cmd.Flags().GetString("llm-model")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			mcpServerFlags, err := cmd.Flags().GetStringArray("mcp-server")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			var cfg *ollamaagent.Config

			if configPath != "" {
				cfg = mustutils.Must(ollamaagent.LoadConfig(ctx, configPath))
			}

			// CLI flags override or provide config
			if ollamaUrl != "" || llmModel != "" || len(mcpServerFlags) > 0 {
				if cfg == nil {
					cfg = &ollamaagent.Config{}
				}

				if ollamaUrl != "" {
					cfg.Ollama.URL = ollamaUrl
				}

				if llmModel != "" {
					cfg.Ollama.Model = llmModel
				}

				if len(mcpServerFlags) > 0 {
					cfg.McpServers = []ollamaagent.McpServerConfig{}
					for _, entry := range mcpServerFlags {
						parts := strings.SplitN(entry, "=", 2)
						if len(parts) != 2 {
							logging.LogFatalf("Invalid --mcp-server format '%s'. Expected 'name=url'.", entry)
						}
						cfg.McpServers = append(cfg.McpServers, ollamaagent.McpServerConfig{
							Name: parts[0],
							URL:  parts[1],
						})
					}
				}
			}

			if cfg == nil {
				logging.LogFatal("Please specify either --config or --ollama-url, --llm-model, and --mcp-server flags.")
			}

			if cfg.Ollama.URL == "" {
				logging.LogFatal("Ollama URL is not set. Use --ollama-url or set it in the config file.")
			}

			if cfg.Ollama.Model == "" {
				logging.LogFatal("LLM model is not set. Use --llm-model or set it in the config file.")
			}

			if len(cfg.McpServers) == 0 {
				logging.LogFatal("No MCP servers configured. Use --mcp-server or set them in the config file.")
			}

			mustutils.Must0(ollamaagent.RunOllamaAgentWithConfig(ctx, cfg))

			logging.LogGoodByCtxf(ctx, "Ollama MCP agent finished.")
		},
	}

	cmd.Flags().String("config", "", "Path to the YAML config file defining the ollama model and MCP servers.")
	cmd.Flags().String("ollama-url", "", "Full URL to the Ollama instance. E.g. http://localhost:11434")
	cmd.Flags().String("llm-model", "", "Name of the LLM model to use. E.g. llama3.1 . For CPU only machines use qwen2.5:3b .")
	cmd.Flags().StringArray("mcp-server", []string{}, "MCP server in the format 'name=url'. Can be specified multiple times.")

	return cmd
}
