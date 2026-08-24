package homeautomationcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/homeautomationcmd/shellycmd"
	"os"
)

func NewHomeAutomationCmd() *cobra.Command {
	const short = "Home automation related commands"

	cmd := &cobra.Command{
		Use:   "homeautomation",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` homeautomation homeautomation`,
	}

	cmd.AddCommand(
		shellycmd.NewShellyCmd(),
	)

	return cmd
}
