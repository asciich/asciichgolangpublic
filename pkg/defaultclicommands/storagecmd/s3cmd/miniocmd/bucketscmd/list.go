package bucketscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func NewListBucketsCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	const short = "List buckets."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage s3 minio buckets list`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			client := options.GetClient(ctx, cmd)

			for _, bucketName := range mustutils.Must(nativeminioclient.ListBucketNames(ctx, client)) {
				fmt.Println(bucketName)
			}

			logging.LogInfoByCtxf(ctx, "List minio buckets finished.")
		},
	}

	return cmd
}
