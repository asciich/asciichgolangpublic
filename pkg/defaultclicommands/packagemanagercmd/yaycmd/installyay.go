package yaycmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
)

func NewInstallYayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-yay",
		Short: "Install the yay package manager itself.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			useSudo, err := cmd.Flags().GetBool("use-sudo")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			mustutils.Must0(yay.InstallYay(ctx, commandexecutorexecoo.Exec(), &packagemanageroptions.InstallPackageOptions{
				UseSudo: useSudo,
			}))
			logging.LogGoodByCtxf(ctx, "Yay installed.")
		},
	}

	cmd.Flags().Bool("use-sudo", false, "Use sudo to get access to install yay.")

	return cmd
}
