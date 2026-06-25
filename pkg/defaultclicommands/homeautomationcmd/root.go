package homeautomationcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/homeautomationcmd/shellycmd"
)

func NewHomeAutomationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "homeautomation",
		Short: "Home automation related commands",
	}

	cmd.AddCommand(
		shellycmd.NewShellyCmd(),
	)

	return cmd
}
