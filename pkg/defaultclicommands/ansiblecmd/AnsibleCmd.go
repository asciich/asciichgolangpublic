package ansiblecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/ansiblecmd/collectioncmd"
	"os"
)

func NewAnsibleCmd() *cobra.Command {
	const short = "Ansible related commands."

	cmd := &cobra.Command{
		Use:   "ansible",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ansible ansible`,
	}

	cmd.AddCommand(
		NewRunRoleCmd(),

		collectioncmd.NewCollectionCmd(),
	)

	return cmd
}
