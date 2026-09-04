package kvmcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmnetworkscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmstorageppolscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmvmscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmvolumescmd"
)

func NewKvmCmd() *cobra.Command {
	const short = "kvm (kernel based virtual machine) related commands."

	cmd := &cobra.Command{
		Use:   "kvm",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvm`,
	}

	cmd.AddCommand(
		kvmnetworkscmd.NewNetworksCmd(),
		kvmstorageppolscmd.NewStoragePoolsCmd(),
		kvmvmscmd.NewVmsCmd(),
		kvmvolumescmd.NewVolumesCmd(),
	)

	return cmd
}
