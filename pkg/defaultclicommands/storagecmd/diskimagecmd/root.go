package diskimagecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewDiskImageCmd() *cobra.Command {
	const short = "Handle disk images like the one created by `dd` or the images to write to SD cards for embedded systems."

	cmd := &cobra.Command{
		Use:   "disk-image",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage diskimage disk-image`,
	}

	cmd.AddCommand(
		NewListPartitionsCmd(),
	)

	return cmd
}
