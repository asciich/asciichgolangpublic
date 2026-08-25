package loggingexamplescmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewLoggingExamplesCmd() *cobra.Command {
	const short = "Examples to showcase the logging functionality."

	cmd := &cobra.Command{
		Use:   "examples",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` logging loggingexamples examples`,
	}

	cmd.AddCommand(
		NewLogInfoCmd(),
		NewLogInfoMultilineCmd(),
	)

	return cmd
}
