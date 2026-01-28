package networkcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/publicipscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/tcpcmd"
)

func NewNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Network related commands.",
	}

	cmd.AddCommand(
		publicipscmd.NewPublicIpsCmd(),
		tcpcmd.NewTcpCmd(),
	)

	return cmd
}
