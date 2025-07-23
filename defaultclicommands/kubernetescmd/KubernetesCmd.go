package kubernetescmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/defaultclicommands/kubernetescmd/eventscmd"
)

func NewKubernetesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "Kubernetes related commands.",
	}

	cmd.AddCommand(
		eventscmd.NewEventsCmd(),
		ListKindNamesCmd(),
	)

	return cmd
}
