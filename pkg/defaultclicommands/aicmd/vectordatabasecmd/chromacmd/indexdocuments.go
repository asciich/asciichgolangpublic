package chromacmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/vectordatabaseutils/chromautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewIndexDocumentsCmd() *cobra.Command {
	const short = "Index documents into the chroma vector database."
	cmd := &cobra.Command{
		Use:   "index-documents",
		Short: short,
		Long: short + `

Example usage for indexing all documents in the current directory:
  ` + os.Args[0] + ` ai vectordatabase chroma index-documents --documents-dir='.' --ollama-url='http://localhost:11434' --chroma-url='http://chroma.example.com' --chroma-collection-name="examplecollection" --basename-regex='.*\\.md$' --verbose
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			documentsDir, err := cmd.Flags().GetString("documents-dir")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if documentsDir == "" {
				logging.LogFatal("Please specify --documents-dir")
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

			if chromaUrl == "" {
				logging.LogFatal("Please specify --chroma-collection-name")
			}

			basenameRegex, err := cmd.Flags().GetString("basename-regex")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if basenameRegex == "" {
				logging.LogFatal("Please specify --basename-regex")
			}

			mustutils.Must0(chromautils.IndexDocuments(
				ctx,
				&chromautils.IndexOptions{
					DocumentsDirectory:   documentsDir,
					OllamaUrl:            ollamaUrl,
					ChromaUrl:            chromaUrl,
					ChromaCollectionName: chromaCollectionName,
					BaseNameRegex:        basenameRegex,
				}))

			logging.LogGoodByCtxf(ctx, "Index documents finished.")
		},
	}

	cmd.Flags().String("documents-dir", "", "Directory containing the documents to collect. Subdirectories are included as well. Use '.' for the current working directory")
	cmd.Flags().String("ollama-url", "", "Full URL to the Ollama to use. E.g. http://localhost:11434")
	cmd.Flags().String("chroma-url", "", "Full URL to the Chroma vector database. e.g http://chroma.example.com")
	cmd.Flags().String("chroma-collection-name", "", "Name of the collection in Chroma to store the vectors.")
	cmd.Flags().String("basename-regex", "", "Regex applied to all filenames in --document-dir. Matching files are indexed, all others are ignored. Example for all MarkDown files: '.*\\.md$'")

	return cmd
}
