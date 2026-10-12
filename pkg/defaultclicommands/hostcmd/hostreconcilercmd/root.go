package hostreconcilercmd

import (
	"github.com/spf13/cobra"
)

func NewHostReconcilerCmd() *cobra.Command {
	const short = "Reconcile the current host according to a given config"

	cmd := &cobra.Command{
		Use:   "hostreconciler",
		Short: short,
		Long:  short + ``,
	}

	cmd.AddCommand(
		NewReconcileOnceCmd(),
	)

	return cmd
}
