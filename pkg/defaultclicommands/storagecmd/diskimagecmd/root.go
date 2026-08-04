package diskimagecmd

import "github.com/spf13/cobra"

func NewDiskImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk-image",
		Short: "Handle disk images like the the one created by `dd` or the images to write to SD cards for embedded systems.",
	}

	cmd.AddCommand(
		NewDiskImageCmd(),
	)

	return cmd
}
