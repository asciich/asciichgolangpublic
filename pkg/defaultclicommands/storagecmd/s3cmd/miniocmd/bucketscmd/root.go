package bucketscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
)

func NewBucketsCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buckets",
		Short: "Buckets related commands",
	}

	cmd.AddCommand(
		NewListBucketsCmd(options),
	)

	return cmd
}
