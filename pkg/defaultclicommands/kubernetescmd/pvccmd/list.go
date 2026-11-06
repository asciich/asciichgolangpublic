package pvccmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListCmd() *cobra.Command {
	const short = "List all kubernetes persistent volume claims."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes pvc list [--namespace <namespace>]

Usage example:
    ` + os.Args[0] + ` kubernetes pvc list
    ` + os.Args[0] + ` kubernetes pvc list --namespace longhorn-system
`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			namespace := mustutils.Must(cmd.Flags().GetString("namespace"))

			cluster := mustutils.Must(nativekubernetesoo.GetDefaultCluster(ctx))

			for _, pvcName := range mustutils.Must(cluster.ListPersistentVolumeClaimNames(ctx, namespace)) {
				fmt.Println(pvcName)
			}
		},
	}

	cmd.Flags().String("namespace", "", "Only list persistent volume claims in --namespace. If not set all namespaces are listed.")

	return cmd
}
