package tracederrorscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func NewDemoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Demonstrate the output of a TracedError.",
		Run: func(cmd *cobra.Command, args []string) {
			cliDemo()
		},
	}

	return cmd
}

func cliDemo() {
	logging.LogGoErrorFatal(tracederrors.TracedError("Example TracedError"))
}