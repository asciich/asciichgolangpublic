package kvmcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmstorageppolscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmvmscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmvolumescmd"
)

func NewKvmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kvm",
		Short: "kvm (kernel based virtual machine) related commands.",
	}

	cmd.AddCommand(
		kvmstorageppolscmd.NewStoragePoolsCmd(),
		kvmvmscmd.NewVmsCmd(),
		kvmvolumescmd.NewVolumesCmd(),
	)

	return cmd
}
