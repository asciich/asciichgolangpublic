package gitlabcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitlabutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewDownloadMainReadmesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-main-readmes",
		Short: "Download the main readmes of the given gitlab group URL.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one gitlab group URL")
			}

			outputDir, err := cmd.Flags().GetString("output-dir")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if outputDir == "" {
				logging.LogFatal("Please specify --output-dir")
			}

			ignoreNoReadmeMd, err := cmd.Flags().GetBool("ignore-no-readme")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			groupUrl := args[0]

			mustutils.Must0(
				gitlabutils.DownloadMainReadmes(ctx, &gitlabutils.DownloadMainReadmesOptions{
					OuputPath:        outputDir,
					GitlabGroupUrl:   groupUrl,
					IgnoreNoReadmeMd: ignoreNoReadmeMd,
				}),
			)

			logging.LogGoodByCtxf(ctx, "Download main README.md files from gitlab group '%s' finished.", groupUrl)
		},
	}

	cmd.PersistentFlags().String("output-dir", "", "Output directory to write the downloaded README.md files in.")
	cmd.PersistentFlags().Bool("ignore-no-readme", false, "Ignore repositories without a README.md.")

	return cmd
}
