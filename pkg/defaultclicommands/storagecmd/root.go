package storagecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/diskimagecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd"
	"os"
)

func NewStorageCmd() *cobra.Command {
	const short = "Storage related commands"

	cmd := &cobra.Command{
		Use:   "storage",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage storage`,
	}

	cmd.AddCommand(
		diskimagecmd.NewDiskImageCmd(),
		s3cmd.NewS3Cmd(),

		NewSpeedTestCmd(),
		NewSyncCmd(),
	)

	return cmd
}
