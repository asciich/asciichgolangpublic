package kindcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewInstallKindCmd() *cobra.Command {
	const short = "Install KinD on the local machine as documented on the official project website https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries ."

	cmd := &cobra.Command{
		Use:   "install",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes kind install [--version <version>]

Usage example:
    ` + os.Args[0] + ` kubernetes kind install --version ` + kindutils.DEFAULT_KIND_VERSION + `
`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			version, err := cmd.Flags().GetString("version")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			mustutils.Must0(kindutils.InstallKind(ctx, version))

			logging.LogGoodByCtxf(ctx, "KinD version '%s' installed.", version)
		},
	}

	cmd.Flags().String("version", kindutils.DEFAULT_KIND_VERSION, "The KinD version to install, e.g. "+kindutils.DEFAULT_KIND_VERSION+".")

	return cmd
}
