package networkcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/dnscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/publicipscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/tcpcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd"
)

func NewNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Network related commands.",
	}

	cmd.AddCommand(
		dnscmd.NewDnsCommand(),
		publicipscmd.NewPublicIpsCmd(),
		tcpcmd.NewTcpCmd(),
		vpncmd.NewVpnCmd(),
	)

	return cmd
}
