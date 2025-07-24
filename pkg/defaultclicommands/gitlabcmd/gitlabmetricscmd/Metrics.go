package gitlabmetricscmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/tracederrors"
)

func NewMetricsCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "metrics",
		Short: "Gitlab metrics related commands",
	}

	cmd.AddCommand(NewExposeMetricsCommand())

	return cmd
}

func AddMetricsCommand(parent *cobra.Command) (err error) {
	if parent == nil {
		return tracederrors.TracedErrorNil("parent")
	}

	parent.AddCommand(NewMetricsCommand())

	return nil
}
