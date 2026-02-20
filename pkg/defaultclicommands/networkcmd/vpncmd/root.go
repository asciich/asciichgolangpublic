package vpncmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/tailscalecmd"
)

func NewVpnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpn",
		Short: "VPN (virtual private network) related commands",
	}

	cmd.AddCommand(
		headscalecmd.NewHeadscaleCmd(),
		tailscalecmd.NewTailscaleCmd(),
	)

	return cmd
}
