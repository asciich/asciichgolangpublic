package kvmvmscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewDeleteCmd() *cobra.Command {
	const short = "Delete a KVM virtual machine."

	cmd := &cobra.Command{
		Use:   "delete [vmName]",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmvms delete my-vm-1
    ` + os.Args[0] + ` virtualmachines kvm kvmvms delete --all`,

		Args: cobra.MaximumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			deleteAll := mustutils.Must(cmd.Flags().GetBool("all"))

			if deleteAll {
				if len(args) > 0 {
					logging.LogFatalf("Do not provide a VM name when using '--all'.")
				}

				vms := mustutils.Must(kvmHypervisor.ListVms(ctx))

				for _, vm := range vms {
					vmName := mustutils.Must(vm.GetCachedName())
					mustutils.Must0(vm.Delete(ctx))
					logging.LogChangedByCtxf(ctx, "Deleted KVM VM '%s'.", vmName)
				}

				logging.LogGoodByCtxf(ctx, "Deleted all '%d' KVM VMs.", len(vms))

				return
			}

			if len(args) != 1 {
				logging.LogFatalf("Provide exactly one VM name, or use '--all' to delete all VMs.")
			}

			vmName := args[0]

			vm := mustutils.Must(kvmHypervisor.GetVmByName(ctx, vmName))
			mustutils.Must0(vm.Delete(ctx))

			logging.LogGoodByCtxf(ctx, "Deleted KVM VM '%s'.", vmName)
		},
	}

	cmd.Flags().Bool("all", false, "Delete all KVM VMs on the hypervisor.")

	return cmd
}
