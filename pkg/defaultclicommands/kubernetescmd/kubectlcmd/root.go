package kubectlcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewKubectlCmd() *cobra.Command {
	const short = "kubectl related commands."

	cmd := &cobra.Command{
		Use:   "kubectl",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes kubectl kubectl`,
	}

	cmd.AddCommand(
		NewInstallKubectlCmd(),
	)

	return cmd
}
