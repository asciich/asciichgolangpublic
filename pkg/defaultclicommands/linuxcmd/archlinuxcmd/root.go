package archlinuxcmd

import "github.com/spf13/cobra"

func NewArchLinuxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archlinux",
		Short: "Archlinux related commands",
	}

	cmd.AddCommand(
		NewIsYayInstalledCmd(),
		NewUpdateArchlinuxKeyringCmd(),
	)

	return cmd
}
