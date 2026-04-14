package kvmvmscmd

import "github.com/spf13/cobra"

func NewVmsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vms",
		Short: "Handle KVM virtual machines.",
	}

	cmd.AddCommand(
		NewListCmd(),
	)

	return cmd
}
