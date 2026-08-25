package tracederrorscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"os"
)

func NewDemoCmd() *cobra.Command {
	const short = "Demonstrate the output of a TracedError."

	cmd := &cobra.Command{
		Use:   "demo",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` errors tracederrors demo`,

		Run: func(cmd *cobra.Command, args []string) {
			cliDemo()
		},
	}

	return cmd
}

func cliDemo() {
	logging.LogGoErrorFatal(tracederrors.TracedError("Example TracedError"))
}
