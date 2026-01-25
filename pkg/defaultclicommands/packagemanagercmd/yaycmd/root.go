package yaycmd

import "github.com/spf13/cobra"

func NewYayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yay",
		Short: "yay (Yet Another Yogurt), a popular AUR helper for Arch Linux related commands.",
	}

	cmd.AddCommand(
		NewInstallYayCmd(),
		NewInstallPackagesCmd(),
		NewRemovePackagesCmd(),
	)

	return cmd
}
