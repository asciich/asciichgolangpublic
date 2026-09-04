package kvmnetworkscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewStartCmd() *cobra.Command {
	const short = "Start a KVM network."

	cmd := &cobra.Command{
		Use:   "start <networkName>",
		Short: short,
		Long: short + `

Starts (activates) the given KVM network. This operation is idempotent:
if the network is already active nothing is changed.

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmnetworks start default`,

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			networkName := args[0]

			mustutils.Must0(kvmHypervisor.StartNetworkByName(ctx, networkName))

			logging.LogGoodByCtxf(ctx, "Started KVM network '%s'.", networkName)
		},
	}

	return cmd
}
