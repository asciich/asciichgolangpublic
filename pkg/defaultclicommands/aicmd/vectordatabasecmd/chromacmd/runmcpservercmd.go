package chromacmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/vectordatabaseutils/chromautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunMcpServerCmd() *cobra.Command {
	const short = "Run the MCP server to expose ChromaDB document search via SSE."
	cmd := &cobra.Command{
		Use:   "run-mcp-server",
		Short: short,
		Long: short + `

Example usage:
  ` + os.Args[0] + ` ai vectordatabase chroma run-mcp-server`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if port <= 0 {
				logging.LogFatal("Please specify a valid --port")
			}

			ollamaUrl, err := cmd.Flags().GetString("ollama-url")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if ollamaUrl == "" {
				logging.LogFatal("Please specify --ollama-url")
			}

			chromaUrl, err := cmd.Flags().GetString("chroma-url")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if chromaUrl == "" {
				logging.LogFatal("Please specify --chroma-url")
			}

			chromaCollectionName, err := cmd.Flags().GetString("chroma-collection-name")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if chromaCollectionName == "" {
				logging.LogFatal("Please specify --chroma-collection-name")
			}

			mustutils.Must0(chromautils.RunDocumentSearchMcpServer(
				ctx,
				&chromautils.McpServerOptions{
					Port:                 port,
					OllamaUrl:            ollamaUrl,
					ChromaUrl:            chromaUrl,
					ChromaCollectionName: chromaCollectionName,
				}))

			logging.LogGoodByCtxf(ctx, "MCP server stopped.")
		},
	}

	cmd.Flags().Int("port", 3001, "Port to run the MCP SSE server on")
	cmd.Flags().String("ollama-url", "", "Full URL to the Ollama to use. E.g. http://localhost:11434")
	cmd.Flags().String("chroma-url", "", "Full URL to the Chroma vector database. e.g http://chroma.example.com")
	cmd.Flags().String("chroma-collection-name", "", "Name of the collection in Chroma to query.")

	return cmd
}
