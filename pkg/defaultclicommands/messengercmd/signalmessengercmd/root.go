package signalmessengercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/messengercmd/signalmessengercmd/signalclirestapicmd"
)

func NewSignalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signal",
		Short: "Signal messenger related commands.",
	}

	cmd.AddCommand(
		signalclirestapicmd.NewSignalCliRestApiCmd(),
	)

	return cmd
}
