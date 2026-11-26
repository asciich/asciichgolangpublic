package gitlabcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabcmdgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabmetricscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabprojectscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/pipelineschedulescmd"
)

func NewGitlabCommand(options gitlabcmdgeneric.GitlabCmdOptions) (cmd *cobra.Command) {
	const short = "Gitlab related commands"

	cmd = &cobra.Command{
		Use:   "gitlab",
		Short: short,
	}

	cmd.AddCommand(
		NewDownloadMainReadmesCmd(),
		NewShowFqdnCmd(options),

		gitlabmetricscmd.NewMetricsCommand(),
		gitlabprojectscmd.NewProjectsCmd(options),
		pipelineschedulescmd.NewPipelineSchedulesCmd(),
	)

	cmd.PersistentFlags().String("url", "", "Url of the Gitlab instance.")

	return cmd
}
