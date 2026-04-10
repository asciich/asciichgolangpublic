package objectscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
)

func NewObjectsCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "objects",
		Short: "S3 objects related commands",
	}

	cmd.AddCommand(
		NewDeleteObjectsCmd(options),
		NewListObjectsCmd(options),
		NewShowDownloadUrlCmd(options),
	)

	cmd.PersistentFlags().String("bucket", "", "Name of the bucket.")

	return cmd
}
