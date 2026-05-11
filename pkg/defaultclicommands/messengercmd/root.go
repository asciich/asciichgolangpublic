package messengercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/messengercmd/signalmessengercmd"
)

func NewMessengerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messenger",
		Short: "Messenger (e.g. Signal) related commands",
	}

	cmd.AddCommand(
		signalmessengercmd.NewSignalCmd(),
	)

	return cmd
}
