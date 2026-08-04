package raspbiancmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/embeddedsystemscmd/raspberrypicmd/raspbiancmd/raspbianimagecmd"
)

func NewRaspbianCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raspbian",
		Short: "Raspbian distribution related commands",
	}

	cmd.AddCommand(
		raspbianimagecmd.NewRaspbianImageCmd(),
	)

	return cmd
}
