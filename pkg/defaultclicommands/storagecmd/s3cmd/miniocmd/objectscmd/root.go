package objectscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
)

func NewObjectsCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "objects",
		Short: "S3 objects related commands",
		Long: `S3 objects related commands.

Usage examples:
  # List all objects in a bucket
  asciich objects list --bucket my-bucket

  # Upload a file to a bucket
  asciich objects upload-file ./local-file.txt object-key --bucket my-bucket

  # Delete objects from a bucket
  asciich objects delete object-key-1 object-key-2 --bucket my-bucket
`,
	}

	cmd.AddCommand(
		NewDeleteObjectsCmd(options),
		NewListObjectsCmd(options),
		NewShowDownloadUrlCmd(options),
		NewUploadFileCmd(options),
	)

	cmd.PersistentFlags().String("bucket", "", "Name of the bucket.")

	return cmd
}
