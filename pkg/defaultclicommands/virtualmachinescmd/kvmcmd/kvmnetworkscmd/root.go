package kvmnetworkscmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewNetworksCmd() *cobra.Command {
	const short = "Handle KVM networks."

	cmd := &cobra.Command{
		Use:   "networks",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmnetworks networks`,
	}

	cmd.AddCommand(
		NewListCmd(),
		NewStartCmd(),
	)

	return cmd
}
