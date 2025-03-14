package prometheuscmd

import "github.com/spf13/cobra"

func NewPrometheusCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "prometheus",
		Short: "Prometheus related commands",
	}

	cmd.AddCommand(NewReadMetricCmd())

	return cmd
}
