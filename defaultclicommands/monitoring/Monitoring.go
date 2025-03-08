package monitoring

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/monitoring/prometheus"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func NewMonitoringCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "monitoring",
		Short: "Monitoring related commands.",
	}

	cmd.AddCommand(prometheus.NewPrometheusCommand())

	return cmd
}

func AddMonitoringGommand(parent *cobra.Command) (err error) {
	if parent == nil {
		return tracederrors.TracedErrorNil("parent")
	}

	parent.AddCommand(NewMonitoringCommand())
	return nil
}
