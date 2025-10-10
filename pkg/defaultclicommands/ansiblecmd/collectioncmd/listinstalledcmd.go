package collectioncmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansiblegalaxyutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListInstalledCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-installed",
		Short: "List installed ansible galaxy collections",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)
			installed := mustutils.Must(ansiblegalaxyutils.ListInstalledCollections(ctx, &ansiblegalaxyutils.ListInstalledCollectionsOptions{
				AnsibleVirtualenvPath: mustutils.Must(cmd.Flags().GetString("virtualenv-path")),
			}))

			for name, version := range installed {
				fmt.Printf("%s: %s\n", name, version)
			}

			logging.LogGoodByCtxf(ctx, "Listed '%d' installed ansible collections.", len(installed))
		},
	}

	cmd.PersistentFlags().String(
		"virtualenv-path", "", "Path to the root directory of the virtualenv containing ansible.",
	)

	return cmd
}
