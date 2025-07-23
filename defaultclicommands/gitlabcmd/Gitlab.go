package gitlabcmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/defaultclicommands/gitlabcmd/gitlabmetricscmd"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/defaultclicommands/gitlabcmd/pipelineschedulescmd"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func NewGitlabCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "gitlab",
		Short: "Gitlab related commands",
	}

	cmd.AddCommand(gitlabmetricscmd.NewMetricsCommand())
	cmd.AddCommand(pipelineschedulescmd.NewPipelineSchedulesCmd())

	return cmd
}

func AddGitlabCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewGitlabCommand())

	return nil
}
