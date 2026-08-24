package vpncmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/tailscalecmd"
	"os"
)

func NewVpnCmd() *cobra.Command {
	const short = "VPN (virtual private network) related commands"

	cmd := &cobra.Command{
		Use:   "vpn",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn vpn`,
	}

	cmd.AddCommand(
		headscalecmd.NewHeadscaleCmd(),
		tailscalecmd.NewTailscaleCmd(),
	)

	return cmd
}
