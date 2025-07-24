package pipelineschedulescmd

import "github.com/spf13/cobra"

func NewPipelineSchedulesCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "pipeline-schedules",
		Short: "Scheduled pipeline related commands.",
	}

	cmd.AddCommand(NewListCommand())

	return cmd
}
