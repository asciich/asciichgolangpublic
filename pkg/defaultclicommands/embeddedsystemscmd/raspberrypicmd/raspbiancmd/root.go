package raspbiancmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/embeddedsystemscmd/raspberrypicmd/raspbiancmd/raspbianimagecmd"
	"os"
)

func NewRaspbianCmd() *cobra.Command {
	const short = "Raspbian distribution related commands"

	cmd := &cobra.Command{
		Use:   "raspbian",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` embeddedsystems raspberrypi raspbian raspbian`,
	}

	cmd.AddCommand(
		raspbianimagecmd.NewRaspbianImageCmd(),
	)

	return cmd
}
