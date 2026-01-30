package exoscaleiamcmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/defaultclicommands/cloudcmd/exoscalecmd/exoscaleiamcmd/exoscaleuserscmd"
)

func NewIamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iam",
		Short: "Exoscale identity and access management.",
	}

	cmd.AddCommand(
		exoscaleuserscmd.NewUsersCmd(),
	)

	return cmd
}
