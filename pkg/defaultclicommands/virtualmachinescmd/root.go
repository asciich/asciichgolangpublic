package virtualmachinescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd"
)

func NewVirtualMachinesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "virtual-machines",
		Short: "Virtual machines and hypervisor related commands",
	}

	cmd.AddCommand(
		kvmcmd.NewKvmCmd(),
	)

	cmd.PersistentFlags().String("hostname", "", "Hostname of the KVM hypervisor. Use 'localhost' to run commands against KVM running on the local machine.")

	return cmd
}
