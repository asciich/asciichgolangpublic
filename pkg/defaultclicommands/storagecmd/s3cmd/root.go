package s3cmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd"
)

func NewS3Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3",
		Short: "S3 (Simple Storage Service) related commands.",
	}

	cmd.AddCommand(
		miniocmd.NewMinioCmd(nil),
	)

	return cmd
}
