package historycmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/bashutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewEnableEmmediateWriteCmd() *cobra.Command {
	const short = "Enables immediate read and write of bash history for the current user. A new bash needs to be started after running this command."

	cmd := &cobra.Command{
		Use:   "enable-immediate-write",
		Short: short,
		Long: short + `

This functionality is useful to sync the history between multiple bash prompts in real time.
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(bashutils.EnableImmediateHistoryReadAndWriteForCurrentUser(ctx))

			logging.LogGoodByCtxf(ctx, "Enabled immediate read and write for bash.")
		},
	}

	return cmd
}
