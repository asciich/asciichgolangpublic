package loggingexamplescmd

import "github.com/spf13/cobra"

func NewLoggingExamplesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "examples",
		Short: "Examples to showcase the logging functionality.",
	}

	cmd.AddCommand(
		NewLogInfoCmd(),
		NewLogInfoMultilineCmd(),
	)

	return cmd
}
