package virtualmachinescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd"
	"os"
)

func NewVirtualMachinesCmd() *cobra.Command {
	const short = "Virtual machines and hypervisor related commands"

	cmd := &cobra.Command{
		Use:   "virtual-machines",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines virtual-machines`,
	}

	cmd.AddCommand(
		kvmcmd.NewKvmCmd(),
	)

	cmd.PersistentFlags().String("hostname", "", "Hostname of the KVM hypervisor. Use 'localhost' to run commands against KVM running on the local machine.")

	return cmd
}
