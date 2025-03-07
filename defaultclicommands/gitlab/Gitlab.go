package gitlab

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/gitlab/metrics"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/gitlab/pipelineschedules"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func NewGitlabCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "gitlab",
		Short: "Gitlab related commands",
	}

	cmd.AddCommand(metrics.NewMetricsCommand())
	cmd.AddCommand(pipelineschedules.NewPipelineSchedulesCmd())

	return cmd
}

func AddGitlabCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewGitlabCommand())

	return nil
}
