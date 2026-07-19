package chromacmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/vectordatabaseutils/chromautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewQueryDocumentsCmd() *cobra.Command {
	const short = "Query Ollama which uses the chroma vector database to query the documents."

	cmd := &cobra.Command{
		Use:   "query-documents",
		Short: short,
		Long: short + `

Example usage for querying documents:
  ` + os.Args[0] + ` ai vectordatabase chroma query-documents --question='How do I create a temporary file?' --ollama-url='http://localhost:11434' --chroma-url='http://chroma.example.com' --chroma-collection-name="examplecollection" --verbose
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			question, err := cmd.Flags().GetString("question")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if question == "" {
				logging.LogFatal("Please specify --question")
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

			llmModel, err := cmd.Flags().GetString("llm-model")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			embeddingModel, err := cmd.Flags().GetString("embedding-model")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			nResults, err := cmd.Flags().GetInt("n-results")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			result := mustutils.Must(chromautils.QueryDocuments(
				ctx,
				&chromautils.QueryOptions{
					OllamaUrl:      ollamaUrl,
					ChromaUrl:      chromaUrl,
					CollectionName: chromaCollectionName,
					EmbeddingModel: embeddingModel,
					LlmModel:       llmModel,
					Question:       question,
					NResults:       nResults,
				}))

			fmt.Printf("\n--- Answer ---\n%s\n", result.Answer)
			fmt.Printf("\n--- Sources ---\n")
			for i, source := range result.SourceFiles {
				distance := float32(0)
				if i < len(result.Distances) {
					distance = result.Distances[i]
				}
				fmt.Printf("  [%d] %s (distance: %.4f)\n", i, source, distance)
			}

			logging.LogGoodByCtxf(ctx, "Query documents finished.")
		},
	}

	cmd.Flags().String("question", "", "The question to ask about the indexed documents.")
	cmd.Flags().String("ollama-url", "", "Full URL to the Ollama to use. E.g. http://localhost:11434")
	cmd.Flags().String("chroma-url", "", "Full URL to the Chroma vector database. e.g http://chroma.example.com")
	cmd.Flags().String("chroma-collection-name", "", "Name of the collection in Chroma to query.")
	cmd.Flags().String("llm-model", "", "LLM model to use for answering. Defaults to 'llama3' if not set.")
	cmd.Flags().String("embedding-model", "", "Embedding model to use for query vectorization. Defaults to 'nomic-embed-text' if not set.")
	cmd.Flags().Int("n-results", 4, "Number of closest document chunks to retrieve from Chroma.")

	return cmd
}
