package kvmvmscmd

import (
	"fmt"
	"os"
	"text/tabwriter"

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

			if wide {
				writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(writer, "NAME\tSTATE\tNETWORK\tMAC\tIP")

				for _, vm := range vms {
					vmName := mustutils.Must(vm.GetCachedName())
					vmState := mustutils.Must(vm.GetCachedState())

					networkName := "-"
					macAddress := "-"
					ipAddress := "-"

					if isRunning := mustutils.Must(vm.IsRunning()); isRunning {
						if network, err := vm.GetNetworkName(ctx); err == nil {
							networkName = network
						}

						if mac, err := vm.GetMacAddress(ctx); err == nil {
							macAddress = mac
						}

						// A VM may be running but not (yet) have a DHCP lease.
						if ip, err := vm.GetIpAddress(ctx); err == nil {
							ipAddress = ip
						}
					}

					fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", vmName, vmState, networkName, macAddress, ipAddress)
				}

				mustutils.Must0(writer.Flush())
			} else {
				for _, vm := range vms {
					fmt.Println(mustutils.Must(vm.GetCachedName()))
				}
			}

			logging.LogGoodByCtxf(ctx, "Listed '%d' KVM VMs.", len(vms))
		},
	}

	cmd.Flags().Bool("wide", false, "Show additional columns: state, network, MAC and IP address.")

	return cmd
}
