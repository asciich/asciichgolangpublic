package kvmvmscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewGetXmlCmd() *cobra.Command {
	const short = "Show the XML definition of a KVM virtual machine."

	cmd := &cobra.Command{
		Use:   "get-xml <vmName>",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmvms get-xml my-vm-1`,

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			vmName := args[0]

			vm := mustutils.Must(kvmHypervisor.GetVmByName(ctx, vmName))
			domainXml := mustutils.Must(vm.GetDomainXmlAsString(ctx))

			fmt.Println(domainXml)

			logging.LogGoodByCtxf(ctx, "Dumped XML of KVM VM '%s'.", vmName)
		},
	}

	return cmd
}
