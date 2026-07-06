package kindcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewInstallKindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "install",
		Short: "Install KinD on the local machine as documented on the official project website https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries .",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(kindutils.InstallKind(ctx))

			logging.LogGoodByCtxf(ctx, "KinD installed")
		},
	}

	return cmd
}
