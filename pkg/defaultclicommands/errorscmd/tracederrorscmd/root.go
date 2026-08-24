package tracederrorscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewTracedErrorsCmd() *cobra.Command {
	const short = "TracedErrors related commands"

	cmd := &cobra.Command{
		Use:   "tracederrors",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` errors tracederrors tracederrors`,
	}

	cmd.AddCommand(
		NewDemoCmd(),
	)

	return cmd
}
