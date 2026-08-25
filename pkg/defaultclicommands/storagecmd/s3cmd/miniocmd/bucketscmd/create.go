package bucketscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
	"os"
)

func NewCreateBucketCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	var publicReadable bool

	const short = "Create a new bucket."

	cmd := &cobra.Command{
		Use:   "create",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage s3 minio buckets create`,

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			client := options.GetClient(ctx, cmd)

			bucketName := args[0]

			createOptions := &s3options.CreateBucketOptions{
				PublicReadable: publicReadable,
			}

			mustutils.Must0(nativeminioclient.CreateBucket(ctx, client, bucketName, createOptions))

			logging.LogInfoByCtxf(ctx, "Create bucket '%s' finished.", bucketName)
		},
	}

	cmd.Flags().BoolVar(&publicReadable, "public-readable", false, "Make the bucket public readable")

	return cmd
}
