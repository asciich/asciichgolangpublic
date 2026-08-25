package bucketscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"os"
)

func NewBucketsCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	const short = "Buckets related commands"

	cmd := &cobra.Command{
		Use:   "buckets",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage s3 minio buckets buckets`,
	}

	cmd.AddCommand(
		NewListBucketsCmd(options),
		NewCreateBucketCmd(options),
	)

	return cmd
}
