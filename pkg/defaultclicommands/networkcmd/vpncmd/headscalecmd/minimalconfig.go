package headscalecmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
)

func NewMinimalConfigCmd() *cobra.Command {
	const short = "Show the minimal headscale config"

	cmd := &cobra.Command{
		Use:   "get-minimal-config",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale get-minimal-config`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(headscalegeneric.GetMinimalDockerConfig())
		},
	}

	return cmd
}
