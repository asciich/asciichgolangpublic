package eventscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewEventsCmd() *cobra.Command {
	const short = "Kubernetes events related commands."

	cmd := &cobra.Command{
		Use:   "events",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes events events`,
	}

	cmd.AddCommand(
		NewWatchCmd(),
	)

	return cmd
}
