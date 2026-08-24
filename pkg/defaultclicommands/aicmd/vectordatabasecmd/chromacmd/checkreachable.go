package chromacmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/vectordatabaseutils/chromautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"os"
)

func NewCheckReachableCmd() *cobra.Command {
	const short = "Check if given --chroma-url instance is reachable."

	cmd := &cobra.Command{
		Use:   "check-reachable",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai vectordatabase chroma check-reachable`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			chromaUrl, err := cmd.Flags().GetString("chroma-url")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			if chromaUrl == "" {
				logging.LogFatal("Please specify --chroma-url")
			}

			chromaClient := chromautils.NewClient(chromaUrl)

			mustutils.Must0(chromaClient.CheckReachable(ctx))

			logging.LogGoodByCtxf(ctx, "Chroma at %s is reachable.", chromaUrl)
		},
	}

	cmd.Flags().String("chroma-url", "", "Full URL to the Chroma vector database. e.g http://chroma.example.com")

	return cmd
}
