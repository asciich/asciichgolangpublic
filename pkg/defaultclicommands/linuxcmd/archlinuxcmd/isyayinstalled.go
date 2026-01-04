package archlinuxcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
)

func NewIsYayInstalledCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-yay-installed",
		Short: "Check if yay is installed. Returns 0 if yay is installed, 1 otherwise.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			ce := commandexecutorexecoo.NewExec()
			isInstalled := mustutils.Must(yay.IsInstalled(ctx, ce))

			if isInstalled {
				logging.LogGoodByCtx(ctx, "yay is installed.")
				os.Exit(0)
			}

			logging.LogErrorByCtx(ctx, "yay is not installed.")
			os.Exit(1)
		},
	}

	return cmd
}
