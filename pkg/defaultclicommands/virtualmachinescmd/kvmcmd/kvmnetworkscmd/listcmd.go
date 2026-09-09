package kvmnetworkscmd

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
	const short = "List KVM networks."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmnetworks list
    ` + os.Args[0] + ` virtualmachines kvm kvmnetworks list --wide`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			wide := mustutils.Must(cmd.Flags().GetBool("wide"))

			if wide {
				networks := mustutils.Must(kvmHypervisor.ListNetworks(ctx))

				writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(writer, "NAME\tSTATE\tAUTOSTART\tPERSISTENT\tNATTED\tFORWARD\tGATEWAY\tNETMASK\tDHCP-RANGE")

				for _, network := range networks {
					networkName := mustutils.Must(network.GetName())
					networkState := mustutils.Must(network.GetState())
					networkAutostart := mustutils.Must(network.GetAutostart())
					networkPersistent := mustutils.Must(network.GetPersistent())

					// Details come from 'net-dumpxml' and may be absent
					// (e.g. isolated networks or networks without DHCP), so
					// missing values are shown as "-" instead of failing.
					forwardMode := "-"
					if mode, err := network.GetForwardMode(ctx); err == nil && mode != "" {
						forwardMode = mode
					}

					natted := "-"
					if isNatted, err := network.IsNatted(ctx); err == nil {
						if isNatted {
							natted = "yes"
						} else {
							natted = "no"
						}
					}

					gateway := "-"
					if ip, err := network.GetGatewayIp(ctx); err == nil {
						gateway = ip
					}

					netmask := "-"
					if mask, err := network.GetNetmask(ctx); err == nil {
						netmask = mask
					}

					dhcpRange := "-"
					if rangeStart, rangeEnd, err := network.GetDhcpRange(ctx); err == nil {
						dhcpRange = fmt.Sprintf("%s-%s", rangeStart, rangeEnd)
					}

					fmt.Fprintf(
						writer,
						"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
						networkName,
						networkState,
						networkAutostart,
						networkPersistent,
						natted,
						forwardMode,
						gateway,
						netmask,
						dhcpRange,
					)
				}

				mustutils.Must0(writer.Flush())

				logging.LogGoodByCtxf(ctx, "Listed '%d' KVM networks.", len(networks))
			} else {
				networkNames := mustutils.Must(kvmHypervisor.ListNetworkNames(ctx))
				for _, networkName := range networkNames {
					fmt.Println(networkName)
				}

				logging.LogGoodByCtxf(ctx, "Listed '%d' KVM networks.", len(networkNames))
			}
		},
	}

	cmd.Flags().Bool("wide", false, "Show additional columns: state, autostart, persistent, NAT, forward mode, gateway, netmask and DHCP range.")

	return cmd
}
