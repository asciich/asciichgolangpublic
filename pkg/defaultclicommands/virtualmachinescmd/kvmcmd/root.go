package kvmcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmstorageppolscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmvmscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmvolumescmd"
	"os"
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
		kvmstorageppolscmd.NewStoragePoolsCmd(),
		kvmvmscmd.NewVmsCmd(),
		kvmvolumescmd.NewVolumesCmd(),
	)

	return cmd
}
