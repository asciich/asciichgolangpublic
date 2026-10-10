package longhorncmd

import (
	"github.com/spf13/cobra"
)

func NewLonghornCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "longhorn",
		Short: "Manage Longhorn distributed block storage for Kubernetes",
		Long: `Commands for Longhorn, the distributed block storage system that provides
persistent volumes (PVCs) for Kubernetes clusters.

Includes tooling such as installing the longhornctl CLI, which is used for
preflight checks and troubleshooting of Longhorn nodes.`,
	}

	cmd.AddCommand(
		NewInstallCliCmd(),
	)

	return cmd
}
