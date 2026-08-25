package linuxcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/linuxcmd/archlinuxcmd"
	"os"
)

func NewLinuxCmd() *cobra.Command {
	const short = "linux related commands."

	cmd := &cobra.Command{
		Use:   "linux",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` linux linux`,
	}

	cmd.AddCommand(
		archlinuxcmd.NewArchLinuxCmd(),
	)

	return cmd
}
