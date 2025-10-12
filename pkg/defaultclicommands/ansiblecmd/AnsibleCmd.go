package ansiblecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/ansiblecmd/collectioncmd"
)

func NewAnsibleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ansible",
		Short: "Ansible related commands.",
	}

	cmd.AddCommand(
		NewRunRoleCmd(),

		collectioncmd.NewCollectionCmd(),
	)

	return cmd
}
