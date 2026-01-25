package yaycmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
)

func NewRemovePackagesCmd() *cobra.Command {
	const short = "Remove packages using yay."

	cmd := &cobra.Command{
		Use:   "remove-packages",
		Short: short,
		Long: short + `

Usage if yay already installed:
    ` + os.Args[0] + `packagemanager yay remove-packages package1 [package2...]

If you want to perform a yay installation first in case yay is not available yet.
This command is idempotent: If yay is already installed it will be used without reinstallation.
	` + os.Args[0] + `packagemanager yay remove-packages --install-yay --use-sudo package1 [package2...]`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)
			ctx = contextutils.WithChangeIndicator(ctx)

			changedExitCode, err := cmd.Flags().GetInt("changed-exit-code")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			installYay, err := cmd.Flags().GetBool("install-yay")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			useSudo, err := cmd.Flags().GetBool("use-sudo")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if len(args) == 0 {
				logging.LogFatalf("Please specify packages to remove.")
			}
			packages := args

			commandExecutor := commandexecutorexecoo.Exec()

			hostDescription, err := commandExecutor.GetHostDescription()
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if installYay {
				mustutils.Must0(yay.InstallYay(ctx, commandExecutor, &packagemanageroptions.InstallPackageOptions{
					UpdateDatabaseFirst: true,
					UseSudo:             useSudo,
				}))
			}

			mustutils.Must0(yay.RemovePackages(ctx, commandExecutor, packages, &packagemanageroptions.RemovePackageOptions{
				UseSudo:             useSudo,
			}))

			logging.LogGoodByCtxf(ctx, "Removeed yay packages '%v' on '%s'.", packages, hostDescription)

			if changedExitCode != 0 {
				osutils.ExitWithChangedExitCode(ctx, changedExitCode)
			}
		},
	}

	cmd.Flags().Int("changed-exit-code", 0, "ExitCode to use in case any change was performed by this command.")
	cmd.Flags().Bool("install-yay", false, "If set, yay will be installed first if not yet available.")
	cmd.Flags().Bool("use-sudo", false, "Use sudo to gain permissions to update when needed during the removeation.")

	return cmd
}
