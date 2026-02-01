package exoscalecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd/exoscalednscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd/exoscaleiamcmd"
)

func NewExoscaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exoscale",
		Short: "Commands for the Exoscale public cloud",
	}

	cmd.AddCommand(
		exoscalednscmd.NewDnsCmd(),
		exoscaleiamcmd.NewIamCmd(),
	)

	return cmd
}
