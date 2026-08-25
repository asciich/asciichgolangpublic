package collectioncmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewCollectionCmd() *cobra.Command {
	const short = "Ansible collection related commands."

	cmd := &cobra.Command{
		Use:   "collection",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ansible collection collection`,
	}

	cmd.AddCommand(
		NewCreateFileStructureCmd(),
		NewListInstalledCmd(),
	)

	return cmd
}
