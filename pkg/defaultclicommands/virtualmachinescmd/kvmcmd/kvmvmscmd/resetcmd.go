package kvmvmscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewResetCmd() *cobra.Command {
	const short = "Reset a KVM virtual machine."

	cmd := &cobra.Command{
		Use:   "reset <vmName>",
		Short: short,
		Long: short + `

Performs a hard reset (like the physical reset button) of the given VM.
This operation is idempotent: if the VM is not running nothing is changed.

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmvms reset <vmname>`,

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			vmName := args[0]

			mustutils.Must0(kvmHypervisor.ResetVm(ctx, vmName))

			logging.LogGoodByCtxf(ctx, "Reset KVM VM '%s'.", vmName)
		},
	}

	return cmd
}
