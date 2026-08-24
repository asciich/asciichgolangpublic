package raspbianimagecmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewRaspbianImageCmd() *cobra.Command {
	const short = "Raspbian image related commands"

	cmd := &cobra.Command{
		Use:   "raspbian-image",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` embeddedsystems raspberrypi raspbian raspbianimage raspbian-image`,
	}

	cmd.AddCommand(
		NewDownloadCmd(),
	)

	return cmd
}
