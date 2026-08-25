package kubernetescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd/eventscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd/k9scmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd/kindcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd/kubectlcmd"
	"os"
)

func NewKubernetesCmd() *cobra.Command {
	const short = "Kubernetes related commands."

	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` kubernetes kubernetes`,
	}

	cmd.AddCommand(
		eventscmd.NewEventsCmd(),
		k9scmd.NewK9sCmd(),
		kindcmd.NewKindCmd(),
		kubectlcmd.NewKubectlCmd(),
		ListKindNamesCmd(),
	)

	return cmd
}
