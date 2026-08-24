package headscalecmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/localdevcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscaleoo"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func NewHeadscaleCmd() *cobra.Command {
	const short = "Headscale is an open source, self-hosted implementation of the Tailscale control server."

	cmd := &cobra.Command{
		Use:   "headscale",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale headscale`,
	}

	cmd.AddCommand(
		localdevcmd.NewLocalDevCmd(),
		operateheadscalecmd.NewOperateCmd(
			&operateheadscalecmd.OperateOptions{
				GetHeadScale: func(ctx context.Context, cmd *cobra.Command) headscaleinterfaces.HeadScale {
					headscale, err := commandexecutorheadscaleoo.NewOnLocalhost()
					if err != nil {
						logging.LogGoErrorFatalWithTrace(err)
					}

					return headscale
				},
			},
		),
		NewMinimalConfigCmd(),
	)

	return cmd
}
