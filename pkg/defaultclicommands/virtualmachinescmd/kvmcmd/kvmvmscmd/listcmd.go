package kvmvmscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListCmd() *cobra.Command {
	const short = "List KVM virtual machines."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmvms list
    ` + os.Args[0] + ` virtualmachines kvm kvmvms list --wide`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			wide := mustutils.Must(cmd.Flags().GetBool("wide"))

			vms := mustutils.Must(kvmHypervisor.ListVms(ctx))
			for _, vm := range vms {
				vmName := mustutils.Must(vm.GetCachedName())

				if wide {
					vmState := mustutils.Must(vm.GetCachedState())
					fmt.Printf("%s\t%s\n", vmName, vmState)
				} else {
					fmt.Println(vmName)
				}
			}

			logging.LogGoodByCtxf(ctx, "Listed '%d' KVM VMs.", len(vms))
		},
	}

	cmd.Flags().Bool("wide", false, "Show additional columns like the VM state.")

	return cmd
}
