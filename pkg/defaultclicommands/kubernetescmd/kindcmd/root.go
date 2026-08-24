package kindcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewKindCmd() *cobra.Command {
	const short = "KinD (Kubernetes in Docker) realted commands."

	cmd := &cobra.Command{
		Use:   "kind",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes kind kind`,
	}

	cmd.AddCommand(
		NewInstallKindCmd(),
	)

	return cmd
}
