package localdevcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscaleoo"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalelocaldevserver"
)

const CONTAINER_NAME = headscalelocaldevserver.DEFAULT_CONTAINER_NAME
const DEFAULT_PORT = headscalelocaldevserver.DEFAULT_PORT

func NewLocalDevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local-dev",
		Short: "Local headscale development environment based on the headscale/headscale docker container.",
	}

	cmd.AddCommand(
		NewRunServerCmd(),
		operateheadscalecmd.NewOperateCmd(
			&operateheadscalecmd.OperateOptions{
				RootCmdShort: "Operate the local development environment.",
				GetHeadScale: func(cmd *cobra.Command) headscaleinterfaces.HeadScale {
					containerName, err := cmd.Flags().GetString("container-name")
					if err != nil {
						logging.LogGoErrorFatalWithTrace(err)
					}

					container, err := nativedocker.GetContainerByName(containerName)
					if err != nil {
						logging.LogGoErrorFatal(err)
					}

					headscale, err := commandexecutorheadscaleoo.New(container)
					if err != nil {
						logging.LogGoErrorFatal(err)
					}

					return headscale
				},
			},
		),
	)

	cmd.PersistentFlags().String("container-name", CONTAINER_NAME, "Name of the docker container running the local development headscale.")
	cmd.PersistentFlags().Int("port", DEFAULT_PORT, "Port for the local headscale dev server.")

	return cmd
}

func getContainerNameAndPort(cmd *cobra.Command) (string, int) {
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		logging.LogGoErrorFatalWithTrace(err)
	}

	containerName, err := cmd.Flags().GetString("container-name")
	if err != nil {
		logging.LogGoErrorFatalWithTrace(err)
	}

	return containerName, port
}
