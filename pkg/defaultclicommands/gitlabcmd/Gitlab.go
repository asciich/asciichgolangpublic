package gitlabcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabmetricscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/pipelineschedulescmd"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func NewGitlabCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "gitlab",
		Short: "Gitlab related commands",
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
