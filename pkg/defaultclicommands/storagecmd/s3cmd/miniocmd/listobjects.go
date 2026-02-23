package miniocmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewListObjectsCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-objects",
		Short: "List the objects in the S3 bucket.",
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

			for _, o := range mustutils.Must(nativeminioclient.ListObjectNames(ctx, client, bucketName)) {
				fmt.Println(o)
			}

			logging.LogGoodByCtxf(ctx, "List objects in bucket '%s' finished.", bucketName)
		},
	}

	cmd.Flags().String("bucket", "", "Name of the bucket to list.")

	return cmd
}
