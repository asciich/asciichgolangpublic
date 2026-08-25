package pipelineschedulescmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewPipelineSchedulesCmd() (cmd *cobra.Command) {
	const short = "Scheduled pipeline related commands."

	cmd = &cobra.Command{
		Use:   "pipeline-schedules",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` gitlab pipelineschedules pipeline-schedules`,
	}

	cmd.AddCommand(NewListCommand())

	return cmd
}
