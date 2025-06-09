package kubernetescmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func ListKindNamesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-kind-names",
		Short: "List all known kind names in the default kubernetes cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			cluster := mustutils.Must(nativekubernetes.GetDefaultCluster(ctx))

			for _, kindName := range mustutils.Must(cluster.ListKindNames(ctx)) {
				fmt.Println(kindName)
			}
		},
	}

	return cmd
}
