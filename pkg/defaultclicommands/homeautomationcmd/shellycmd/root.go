package shellycmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/homeautomationcmd/shellycmd/gen3handtcmd"
	"os"
)

func NewShellyCmd() *cobra.Command {
	const short = "Shelly homeautomation devices related commands."

	cmd := &cobra.Command{
		Use:   "shelly",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` homeautomation shelly shelly`,
	}

	cmd.AddCommand(
		gen3handtcmd.NewGen3HAndTCmd(),
	)

	return cmd
}
