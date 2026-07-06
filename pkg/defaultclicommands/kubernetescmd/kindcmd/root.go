package kindcmd

import "github.com/spf13/cobra"

func NewKindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind",
		Short: "KinD (Kubernetes in Docker) realted commands.",
	}

	cmd.AddCommand(
		NewInstallKindCmd(),
	)

	return cmd
}
