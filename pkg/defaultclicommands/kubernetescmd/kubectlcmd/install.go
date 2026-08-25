package kubectlcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"os"
)

func NewInstallKubectlCmd() *cobra.Command {
	const short = "Install kubectl on current system."

	cmd := &cobra.Command{
		Use:   "install",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes kubectl install`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(kubectlutils.InstallKubectl(ctx, kubectlutils.DefaultInstallKubectlOptions()))

			logging.LogGoodByCtxf(ctx, "kubectl installed.")
		},
	}

	return cmd
}
