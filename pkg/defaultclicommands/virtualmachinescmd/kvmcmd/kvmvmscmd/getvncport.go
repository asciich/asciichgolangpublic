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

func NewGetVncPortCmd() *cobra.Command {
	const short = "Show the VNC port of a KVM virtual machine."

	cmd := &cobra.Command{
		Use:   "get-vnc-port [vmName]",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` virtualmachines kvm kvmvms get-vnc-port my-vm-1
    ` + os.Args[0] + ` virtualmachines kvm kvmvms get-vnc-port -l`,

		Args: cobra.MaximumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			listAll := mustutils.Must(cmd.Flags().GetBool("list"))

			if listAll {
				if len(args) > 0 {
					logging.LogFatalf("Do not provide a VM name when using '--list'/'-l'.")
				}

				vms := mustutils.Must(kvmHypervisor.ListVms(ctx))

				writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(writer, "NAME\tVNC_PORT")

				for _, vm := range vms {
					vmName := mustutils.Must(vm.GetCachedName())

					// A VM may have no VNC port (e.g. not running or no VNC graphics device).
					vncPort := "-"
					if port, err := vm.GetVncPort(ctx); err == nil {
						vncPort = fmt.Sprintf("%d", port)
					}

					fmt.Fprintf(writer, "%s\t%s\n", vmName, vncPort)
				}

				mustutils.Must0(writer.Flush())

				logging.LogGoodByCtxf(ctx, "Listed VNC ports of '%d' KVM VMs.", len(vms))

				return
			}

			if len(args) != 1 {
				logging.LogFatalf("Provide exactly one VM name, or use '--list'/'-l' to list all VMs.")
			}

			vmName := args[0]

			vm := mustutils.Must(kvmHypervisor.GetVmByName(ctx, vmName))
			vncPort := mustutils.Must(vm.GetVncPort(ctx))

			fmt.Println(vncPort)

			logging.LogGoodByCtxf(ctx, "Got VNC port '%d' of KVM VM '%s'.", vncPort, vmName)
		},
	}

	cmd.Flags().BoolP("list", "l", false, "List all VMs and their VNC port.")

	return cmd
}
