package kvmvolumescmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd/kvmcmd/kvmcmdutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		Short: "List KVM volumes",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, kvmHypervisor := kvmcmdutils.GetCtxAndKvmHypervisor(cmd)

			volumeNames := mustutils.Must(kvmHypervisor.GetVolumeNames(ctx))

			for _, vn := range volumeNames {
				fmt.Println(vn)
			}

			logging.LogGoodByCtxf(ctx, "List '%d' KVM volume names finished.", len(volumeNames))
		},
	}

	return cmd
}