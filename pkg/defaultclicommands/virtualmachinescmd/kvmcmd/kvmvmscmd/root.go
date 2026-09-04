package kvmvmscmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewVmsCmd() *cobra.Command {
	const short = "Handle KVM virtual machines."

	cmd := &cobra.Command{
		Use:   "vms",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmvms vms`,
	}

	cmd.AddCommand(
		NewGetXmlCmd(),
		NewListCmd(),
		NewResetCmd(),
	)

	return cmd
}
