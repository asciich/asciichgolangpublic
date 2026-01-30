package exoscalecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd/exoscaleiamcmd"
)

func NewExoscaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exoscale",
		Short: "Commands for the Exoscale public cloud",
	}

	cmd.AddCommand(
		exoscaleiamcmd.NewIamCmd(),
	)

	return cmd
}
