package prometheuscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewPrometheusCommand() (cmd *cobra.Command) {
	const short = "Prometheus related commands"

	cmd = &cobra.Command{
		Use:   "prometheus",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` monitoring prometheus prometheus`,
	}

	cmd.AddCommand(NewReadMetricCmd())

	return cmd
}
