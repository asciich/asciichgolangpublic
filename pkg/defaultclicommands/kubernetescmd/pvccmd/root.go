package pvccmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewPvcCmd() *cobra.Command {
	const short = "Kubernetes persistent volume claims related commands."

	cmd := &cobra.Command{
		Use:   "pvc",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes pvc <list> [flags]`,
	}

	cmd.AddCommand(
		NewListCmd(),
		NewListFilesInPvcCmd(),
	)

	return cmd
}
