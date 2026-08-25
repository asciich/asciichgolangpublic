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
    ` + os.Args[0] + ` virtualmachines kvm kvmvms list`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			vmNames := mustutils.Must(kvmHypervisor.ListVmNames(ctx))
			for _, vmName := range vmNames {
				fmt.Println(vmName)
			}

			logging.LogGoodByCtxf(ctx, "Listed '%d' KVM VMs.", len(vmNames))
		},
	}

	return cmd
}
