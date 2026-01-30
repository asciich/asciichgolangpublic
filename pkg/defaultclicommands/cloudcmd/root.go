package cloudcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd/exoscalecmd"
)

func NewCloudCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cloud",
		Short: "(Public-) cloud related commands.",
	}

	cmd.AddCommand(
		exoscalecmd.NewExoscaleCmd(),
	)

	return cmd
}