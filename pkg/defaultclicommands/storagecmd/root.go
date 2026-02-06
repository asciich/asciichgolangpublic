package storagecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd"
)

func NewStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "storage",
		Short: "Storage related commands",
	}

	cmd.AddCommand(
		s3cmd.NewS3Cmd(),
	)

	return cmd
}
