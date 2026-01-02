package linuxcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/linuxcmd/archlinuxcmd"
)

func NewLinuxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "linux",
		Short: "linux related commands.",
	}

	cmd.AddCommand(
		archlinuxcmd.NewArchLinuxCmd(),
	)

	return cmd
}
