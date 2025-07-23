package ansiblecmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/defaultclicommands/ansiblecmd/collectioncmd"
)

func NewAnsibleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ansible",
		Short: "Ansible related commands.",
	}

	cmd.AddCommand(
		collectioncmd.NewCollectionCmd(),
	)

	return cmd
}
