package kubectlcmd

import "github.com/spf13/cobra"

func NewKubectlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl",
		Short: "kubectl related commands.",
	}

	cmd.AddCommand(
		NewInstallKubectlCmd(),
	)

	return cmd
}
