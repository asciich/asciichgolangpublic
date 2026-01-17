package yaycmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
)

func NewInstallYayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-yay",
		Short: "Install the yay package manager itself.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)
			ctx = contextutils.WithChangeIndicator(ctx)

			changedExitCode, err := cmd.Flags().GetInt("changed-exit-code")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			useSudo, err := cmd.Flags().GetBool("use-sudo")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			mustutils.Must0(yay.InstallYay(ctx, commandexecutorexecoo.Exec(), &packagemanageroptions.InstallPackageOptions{
				UseSudo: useSudo,
			}))
			logging.LogGoodByCtxf(ctx, "Yay installed.")

			if changedExitCode != 0 {
				osutils.ExitWithChangedExitCode(ctx, changedExitCode)
			}
		},
	}

	cmd.Flags().Int("changed-exit-code", 0, "ExitCode to use in case any change was performed by this command.")
	cmd.Flags().Bool("use-sudo", false, "Use sudo to get access to install yay.")

	return cmd
}
