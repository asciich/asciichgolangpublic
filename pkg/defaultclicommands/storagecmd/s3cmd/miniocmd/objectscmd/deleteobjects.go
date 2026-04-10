package objectscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewDeleteObjectsCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the given objects in the S3 --bucket.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			client := options.GetClient(ctx, cmd)

			bucketName, err := cmd.Flags().GetString("bucket")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if bucketName == "" {
				logging.LogFatalf("Please specify --bucket.")
			}

			if len(args) <= 0 {
				logging.LogFatalf("Please specify at least one object to delete.")
			}

			for _, objectKey := range args {
				nativeminioclient.DeleteObject(ctx, client, bucketName, objectKey)
			}

			logging.LogGoodByCtxf(ctx, "List objects in bucket '%s' finished.", bucketName)
		},
	}



	return cmd
}
