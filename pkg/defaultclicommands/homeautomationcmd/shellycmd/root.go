package shellycmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/homeautomationcmd/shellycmd/gen3handtcmd"
)

func NewShellyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shelly",
		Short: "Shelly homeautomation devices related commands.",
	}

	cmd.AddCommand(
		gen3handtcmd.NewGen3HAndTCmd(),
	)

	return cmd
}
