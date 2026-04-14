package kvmvolumescmd

import "github.com/spf13/cobra"

func NewVolumesCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use: "volumes",
		Short: "KVM volumes related commands.",
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}