package objectscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewShowDownloadUrlCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-download-url",
		Short: "Show the download URL of the specified object key in the --bucket",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			bucket, err := cmd.Flags().GetString("bucket")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if bucket == "" {
				logging.LogFatal("Please specify --bucket.")
			}

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one object key to build the download URL.")
			}

			client := options.GetClient(ctx, cmd)

			url := mustutils.Must(nativeminioclient.GetDownloadUrl(ctx, client, bucket, args[0]))
			fmt.Println(url)
		},
	}

	return cmd
}
