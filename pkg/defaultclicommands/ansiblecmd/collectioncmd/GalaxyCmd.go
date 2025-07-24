package collectioncmd

import "github.com/spf13/cobra"

func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Ansible collection related commands.",
	}

	cmd.AddCommand(
		NewCreateFileStructureCmd(),
	)

	return cmd
}
