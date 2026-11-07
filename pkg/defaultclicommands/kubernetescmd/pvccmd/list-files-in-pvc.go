package pvccmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListFilesInPvcCmd() *cobra.Command {
	const short = "List all files inside a kubernetes persistent volume claim."

	cmd := &cobra.Command{
		Use:   "list-files-in-pvc",
		Short: short,
		Long: short + `

A temporary pod is started to list all files inside the given persistent volume claim.
Hidden files are included and the output is sorted alphabetically.

Usage:
    ` + os.Args[0] + ` kubernetes pvc list-files-in-pvc --pvc <pvc-name> [--namespace <namespace>]

Usage example:
    ` + os.Args[0] + ` kubernetes pvc list-files-in-pvc --pvc my-pvc --namespace longhorn-system
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			pvcName := mustutils.Must(cmd.Flags().GetString("pvc"))
			namespace := mustutils.Must(cmd.Flags().GetString("namespace"))

			cluster := mustutils.Must(nativekubernetesoo.GetDefaultCluster(ctx))

			for _, fileName := range mustutils.Must(cluster.ListPersistentVolumeClaimFileNames(ctx, namespace, pvcName)) {
				fmt.Println(fileName)
			}
		},
	}

	cmd.Flags().String("pvc", "", "Name of the persistent volume claim to list the files from")
	cmd.Flags().String("namespace", "", "Namespace of the persistent volume claim. If not set the default namespace is used.")

	mustutils.Must0(cmd.MarkFlagRequired("pvc"))

	return cmd
}
