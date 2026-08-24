package s3cmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd"
	"os"
)

func NewS3Cmd() *cobra.Command {
	const short = "S3 (Simple Storage Service) related commands."

	cmd := &cobra.Command{
		Use:   "s3",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage s3 s3`,
	}

	cmd.AddCommand(
		miniocmd.NewMinioCmd(nil),
	)

	return cmd
}
