package raspberrypicmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/embeddedsystemscmd/raspberrypicmd/raspbiancmd"
)

func NewRaspberryPiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raspberry-pi",
		Short: "Raspberry Pi related commands",
	}

	cmd.AddCommand(
		raspbiancmd.NewRaspbianCmd(),
	)

	return cmd
}
