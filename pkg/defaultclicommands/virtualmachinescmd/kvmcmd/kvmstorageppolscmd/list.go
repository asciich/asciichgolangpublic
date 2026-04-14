package kvmstorageppolscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List storage pools",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			storagePoolNames := mustutils.Must(kvmHypervisor.ListStoragePoolNames(ctx))

			for _, spn := range storagePoolNames {
				fmt.Println(spn)
			}

			logging.LogGoodByCtxf(ctx, "Listed '%d' storage pools.", len(storagePoolNames))
		},
	}

	return cmd
}
