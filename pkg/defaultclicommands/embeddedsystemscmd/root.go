package embeddedsystemscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/embeddedsystemscmd/raspberrypicmd"
)

func NewEmbeddedSystemsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "embedded-systems",
		Short: "Embedded systems related commands",
	}

	cmd.AddCommand(
		raspberrypicmd.NewRaspberryPiCmd(),
	)

	return cmd
}
