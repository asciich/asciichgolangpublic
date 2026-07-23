package routercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/routercmd/pfsensecmd"
)

func NewRouterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "router",
		Short: "Router related commands.",
	}

	cmd.AddCommand(
		pfsensecmd.NewPfSenseCmd(),
	)

	return cmd
}