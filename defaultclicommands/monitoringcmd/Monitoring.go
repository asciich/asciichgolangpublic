package monitoringcmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/defaultclicommands/monitoringcmd/prometheuscmd"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func NewMonitoringCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "monitoring",
		Short: "Monitoring related commands.",
	}

	cmd.AddCommand(prometheuscmd.NewPrometheusCommand())

	return cmd
}

func AddMonitoringGommand(parent *cobra.Command) (err error) {
	if parent == nil {
		return tracederrors.TracedErrorNil("parent")
	}

	parent.AddCommand(NewMonitoringCommand())
	return nil
}
