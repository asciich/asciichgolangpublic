package imagescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewDownloadAsArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-as-archive",
		Short: "Download a container image from a registry to a local archive file. This does no require docker or another deamon running.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.ContextVerbose()

			imageAndTag, err := cmd.Flags().GetString("image")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if imageAndTag == "" {
				logging.LogFatal("Please specify --image")
			}

			outputPath, err := cmd.Flags().GetString("output")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if outputPath == "" {
				logging.LogFatal("Please specify --output")
			}

			mustutils.Must0(containerimagehandler.DownloadImageAsArchive(ctx, imageAndTag, outputPath))

			logging.LogGoodByCtxf(ctx, "Downloaded container image '%s' as local archive '%s'.", imageAndTag, outputPath)
		},
	}

	cmd.Flags().String("image", "", "Image and tag to download. Format is '<name>:<tag>'. If only the name is specified the ':latest' tag is taken automatically.")
	cmd.Flags().String("output", "", "Output path to store the iamge.")

	return cmd
}
