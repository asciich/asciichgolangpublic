package messengercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/messengercmd/signalmessengercmd"
	"os"
)

func NewMessengerCmd() *cobra.Command {
	const short = "Messenger (e.g. Signal) related commands"

	cmd := &cobra.Command{
		Use:   "messenger",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` messenger messenger`,
	}

	cmd.AddCommand(
		signalmessengercmd.NewSignalCmd(),
	)

	return cmd
}
