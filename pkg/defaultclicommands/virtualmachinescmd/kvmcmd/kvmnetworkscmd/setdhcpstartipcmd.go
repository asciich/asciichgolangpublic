package kvmnetworkscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewSetDhcpStartIpCmd() *cobra.Command {
	const short = "Set the DHCP range start IP of a KVM network."

	cmd := &cobra.Command{
		Use:   "set-dhcp-start-ip",
		Short: short,
		Long: short + `

The end of the DHCP range stays unchanged. This is useful to free up low
addresses (e.g. for static IPs or a kube-vip VIP) that libvirt should no
longer hand out via DHCP.

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmnetworks set-dhcp-start-ip --name=default --ip=192.168.122.100`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			networkName := mustutils.Must(cmd.Flags().GetString("name"))
			startIp := mustutils.Must(cmd.Flags().GetString("ip"))

			network := mustutils.Must(kvmHypervisor.GetNetworkByName(networkName))

			mustutils.Must0(network.SetDhcpStartIp(ctx, startIp))

			logging.LogGoodByCtxf(ctx, "Set DHCP range start IP of KVM network '%s' to '%s'.", networkName, startIp)
		},
	}

	cmd.Flags().String("name", "", "Name of the KVM network (e.g. 'default').")
	cmd.Flags().String("ip", "", "New DHCP range start IP address (e.g. '192.168.122.100').")

	mustutils.Must0(cmd.MarkFlagRequired("name"))
	mustutils.Must0(cmd.MarkFlagRequired("ip"))

	return cmd
}
