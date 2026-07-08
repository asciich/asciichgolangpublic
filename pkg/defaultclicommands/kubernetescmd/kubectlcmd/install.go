package kubectlcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewInstallKubectlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install kubectl on current system.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(kubectlutils.InstallKubectl(ctx))

			logging.LogGoodByCtxf(ctx, "kubectl installed.")
		},
	}

	return cmd
}
