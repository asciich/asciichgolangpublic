package kubernetescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd/eventscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd/kindcmd"
)

func NewKubernetesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "Kubernetes related commands.",
	}

	cmd.AddCommand(
		eventscmd.NewEventsCmd(),
		kindcmd.NewKindCmd(),
		ListKindNamesCmd(),
	)

	return cmd
}
