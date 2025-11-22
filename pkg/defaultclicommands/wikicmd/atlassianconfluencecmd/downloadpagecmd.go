package atlassianconfluencecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/wikiutils/atlassianconfluenceutils"
)

func NewDownloadPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-page",
		Short: "Download the given wiki page into the --output-dir.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one URL.")
			}

			url := args[0]

			outputDir, err := cmd.Flags().GetString("output-dir")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if outputDir == "" {
				logging.LogFatal("Please specify --output-dir")
			}

			recursive, err := cmd.Flags().GetBool("recursive")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			mustutils.Must0(atlassianconfluenceutils.DownloadPageContent(ctx, url, outputDir, &atlassianconfluenceutils.DownloadPageContentOptions{
				Recursive: recursive,
			}))

			logging.LogGoodByCtxf(ctx, "Download wiki page %s finished.", url)
		},
	}

	cmd.PersistentFlags().String("output-dir", "", "The output directory to write the downloaded wiki page.")
	cmd.PersistentFlags().Bool("recursive", false, "If set the child pages are downloaded as well.")

	return cmd
}
