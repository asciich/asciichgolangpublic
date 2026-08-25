package exoscalecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd/exoscalednscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd/exoscaleiamcmd"
	"os"
)

func NewExoscaleCmd() *cobra.Command {
	const short = "Commands for the Exoscale public cloud"

	cmd := &cobra.Command{
		Use:   "exoscale",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` cloud exoscale exoscale`,
	}

	cmd.AddCommand(
		exoscalednscmd.NewDnsCmd(),
		exoscaleiamcmd.NewIamCmd(),
	)

	return cmd
}
