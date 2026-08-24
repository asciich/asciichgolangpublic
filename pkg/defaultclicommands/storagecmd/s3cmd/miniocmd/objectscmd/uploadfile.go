package objectscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewUploadFileCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	const short = "Upload a local file into the --bucket."

	cmd := &cobra.Command{
		Use:   "upload-file",
		Short: short,
		Long: `Upload a local file into the --bucket.

Usage example:
  # Upload a local file to a bucket with a specific object key
  asciich objects upload-file ./local-file.txt object-key --bucket my-bucket
`,
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

			if len(args) != 2 {
				logging.LogFatalf("Please specify the local file path and the object key.")
			}

			localFilePath := args[0]
			objectKey := args[1]

			err = nativeminioclient.UploadFileByPath(ctx, client, bucketName, objectKey, localFilePath)
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			logging.LogGoodByCtxf(ctx, "Upload of '%s' to bucket '%s' as '%s' finished.", localFilePath, bucketName, objectKey)
		},
	}

	return cmd
}
