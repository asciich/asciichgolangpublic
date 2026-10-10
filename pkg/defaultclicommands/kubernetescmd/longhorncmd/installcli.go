package longhorncmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/longhornutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewInstallCliCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-cli",
		Short: "Install longhornctl CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(longhornutils.InstallLonghornctl(ctx))

			logging.LogGoodByCtxf(ctx, "Longhorn CLI successfully installed.")
		},
	}

	return cmd
}
