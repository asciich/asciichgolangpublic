package userscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewWhoAmICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show current user.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			cluster := mustutils.Must(nativekubernetesoo.GetDefaultCluster(ctx))
			user := mustutils.Must(cluster.WhoAmI(ctx))
			fmt.Println(user.Username)
		},
	}

	return cmd
}
