package k9scmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/k9sutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewInstallK9sCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "install",
		Short: "Install k9s from https://github.com/derailed/k9s/releases .",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(k9sutils.InstallK9s(ctx))

			logging.LogGoodByCtxf(ctx, "Installed k9s")
		},
	}

	return cmd
}