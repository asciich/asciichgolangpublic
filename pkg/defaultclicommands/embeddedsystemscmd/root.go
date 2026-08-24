package embeddedsystemscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/embeddedsystemscmd/raspberrypicmd"
	"os"
)

func NewEmbeddedSystemsCmd() *cobra.Command {
	const short = "Embedded systems related commands"

	cmd := &cobra.Command{
		Use:   "embedded-systems",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` embeddedsystems embedded-systems`,
	}

	cmd.AddCommand(
		raspberrypicmd.NewRaspberryPiCmd(),
	)

	return cmd
}
