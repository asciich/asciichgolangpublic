package headscalecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/localdevcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscaleoo"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func NewHeadscaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "headscale",
		Short: "Headscale is an open source, self-hosted implementation of the Tailscale control server.",
	}

	cmd.AddCommand(
		localdevcmd.NewLocalDevCmd(),
		operateheadscalecmd.NewOperateCmd(
			&operateheadscalecmd.OperateOptions{
				GetHeadScale: func(cmd *cobra.Command) headscaleinterfaces.HeadScale {
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
