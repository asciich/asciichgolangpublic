package gitlabcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabmetricscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/pipelineschedulescmd"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"os"
)

func NewGitlabCommand() (cmd *cobra.Command) {
	const short = "Gitlab related commands"

	cmd = &cobra.Command{
		Use:   "gitlab",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` gitlab gitlab`,
	}

	cmd.AddCommand(
		NewDownloadMainReadmesCmd(),

		gitlabmetricscmd.NewMetricsCommand(),
		pipelineschedulescmd.NewPipelineSchedulesCmd(),
	)

	return cmd
}

func AddGitlabCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(
		NewGitlabCommand(),
	)

	return nil
}
