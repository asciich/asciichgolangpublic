package eventscmd

import "github.com/spf13/cobra"

func NewEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Kubernetes events related commands.",
	}

	cmd.AddCommand(
		NewWatchCmd(),
	)

	return cmd
}
