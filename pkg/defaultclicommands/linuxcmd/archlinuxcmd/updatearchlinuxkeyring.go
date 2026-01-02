package archlinuxcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/linuxutils/archlinuxutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
)

func NewUpdateArchlinuxKeyringCmd() *cobra.Command {
	const short = "Update the 'archlinux-keyring' package."

	runbook := archlinuxutils.NewUpdateArchLinuxKeyringPackageRunbook(nil, false)
	runbookDocumentation, err := runbook.DocumentSteps()
	if err != nil {
		panic(err)
	}

	cmd := &cobra.Command{
		Use:   "update-archlinux-keyring",
		Short: short,
		Long: short + `
Useful to update the archlinux-keyring containing all signing keys before a system update is done to avoid outdated or missing keys.

This is achieved by:
` + runbookDocumentation + `
`,
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

			commandExectuor := commandexecutorexecoo.NewExec()
			runbook := archlinuxutils.NewUpdateArchLinuxKeyringPackageRunbook(commandExectuor, useSudo)
			mustutils.Must0(runbook.Execute(ctx))

			logging.LogGoodByCtx(ctx, "archlinux-keyring updated.")

			if changedExitCode != 0 {
				osutils.ExitWithChangedExitCode(ctx, changedExitCode)
			}
		},
	}

	cmd.Flags().Int("changed-exit-code", 0, "ExitCode to use in case any change was performed by this command.")
	cmd.Flags().Bool("use-sudo", false, "Use sudo to run archlinux-keyring update.")

	return cmd
}
