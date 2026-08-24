package signalmessengercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/messengercmd/signalmessengercmd/signalclirestapicmd"
	"os"
)

func NewSignalCmd() *cobra.Command {
	const short = "Signal messenger related commands."

	cmd := &cobra.Command{
		Use:   "signal",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` messenger signalmessenger signal`,
	}

	cmd.AddCommand(
		signalclirestapicmd.NewSignalCliRestApiCmd(),
	)

	return cmd
}
