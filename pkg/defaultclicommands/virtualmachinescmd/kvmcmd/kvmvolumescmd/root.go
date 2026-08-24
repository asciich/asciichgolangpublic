package kvmvolumescmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewVolumesCmd() *cobra.Command {
	const short = "KVM volumes related commands."

	cmd := &cobra.Command{
		Use:   "volumes",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmvolumes volumes`,
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
