package tracederrorscmd

import "github.com/spf13/cobra"

func NewTracedErrorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tracederrors",
		Short: "TracedErrors related commands",
	}

	cmd.AddCommand(
		NewDemoCmd(),
	)

	return cmd
}
