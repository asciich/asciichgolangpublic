package kubernetescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/kubernetescmd/eventscmd"
)

func NewKubernetesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "Kubernetes related commands.",
	}

	cmd.AddCommand(
		eventscmd.NewEventsCmd(),
	)

	return cmd
}
