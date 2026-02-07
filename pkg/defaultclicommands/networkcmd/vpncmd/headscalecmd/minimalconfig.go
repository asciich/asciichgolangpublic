package headscalecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
)

func NewMinimalConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-minimal-config",
		Short: "Show the minimal headscale config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(headscalegeneric.GetMinimalDockerConfig())
		},
	}

	return cmd
}
