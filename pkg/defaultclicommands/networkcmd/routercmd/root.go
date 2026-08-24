package routercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/routercmd/pfsensecmd"
	"os"
)

func NewRouterCmd() *cobra.Command {
	const short = "Router related commands."

	cmd := &cobra.Command{
		Use:   "router",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network router router`,
	}

	cmd.AddCommand(
		pfsensecmd.NewPfSenseCmd(),
	)

	return cmd
}
