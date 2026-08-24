package cloudcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd"
	"os"
)

func NewCloudCmd() *cobra.Command {
	const short = "(Public-) cloud related commands."

	cmd := &cobra.Command{
		Use:   "cloud",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` cloud cloud`,
	}

	cmd.AddCommand(
		exoscalecmd.NewExoscaleCmd(),
	)

	return cmd
}
