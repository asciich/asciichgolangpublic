package kvmstorageppolscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewStoragePoolsCmd() *cobra.Command {
	const short = "KVM storage pools related commands"

	cmd := &cobra.Command{
		Use:   "storage-pools",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmstorageppols storage-pools`,
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
