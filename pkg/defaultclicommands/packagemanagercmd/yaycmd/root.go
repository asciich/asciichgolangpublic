package yaycmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewYayCmd() *cobra.Command {
	const short = "yay (Yet Another Yogurt), a popular AUR helper for Arch Linux related commands."

	cmd := &cobra.Command{
		Use:   "yay",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` packagemanager yay yay`,
	}

	cmd.AddCommand(
		NewInstallYayCmd(),
		NewInstallPackagesCmd(),
		NewRemovePackagesCmd(),
	)

	return cmd
}
