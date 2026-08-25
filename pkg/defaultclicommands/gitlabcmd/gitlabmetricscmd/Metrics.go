package gitlabmetricscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"os"
)

func NewMetricsCommand() (cmd *cobra.Command) {
	const short = "Gitlab metrics related commands"

	cmd = &cobra.Command{
		Use:   "metrics",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` gitlab gitlabmetrics metrics`,
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
