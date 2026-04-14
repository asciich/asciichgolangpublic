package kvmstorageppolscmd

import "github.com/spf13/cobra"

func NewStoragePoolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage-pools",
		Short: "KVM storage pools related commands",
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
