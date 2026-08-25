package raspberrypicmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/embeddedsystemscmd/raspberrypicmd/raspbiancmd"
	"os"
)

func NewRaspberryPiCmd() *cobra.Command {
	const short = "Raspberry Pi related commands"

	cmd := &cobra.Command{
		Use:   "raspberry-pi",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` embeddedsystems raspberrypi raspberry-pi`,
	}

	cmd.AddCommand(
		raspbiancmd.NewRaspbianCmd(),
	)

	return cmd
}
