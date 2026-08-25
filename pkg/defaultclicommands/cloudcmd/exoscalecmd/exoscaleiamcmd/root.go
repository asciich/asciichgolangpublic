package exoscaleiamcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd/exoscaleiamcmd/exoscaleuserscmd"
	"os"
)

func NewIamCmd() *cobra.Command {
	const short = "Exoscale identity and access management."

	cmd := &cobra.Command{
		Use:   "iam",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` cloud exoscale exoscaleiam iam`,
	}

	cmd.AddCommand(
		exoscaleuserscmd.NewUsersCmd(),
	)

	return cmd
}
