package archlinuxcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewArchLinuxCmd() *cobra.Command {
	const short = "Archlinux related commands"

	cmd := &cobra.Command{
		Use:   "archlinux",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` linux archlinux archlinux`,
	}

	cmd.AddCommand(
		NewIsYayInstalledCmd(),
		NewUpdateArchlinuxKeyringCmd(),
	)

	return cmd
}
