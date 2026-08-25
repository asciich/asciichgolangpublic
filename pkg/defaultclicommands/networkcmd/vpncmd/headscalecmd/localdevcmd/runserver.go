package localdevcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalelocaldevserver"
	"github.com/asciich/asciichgolangpublic/pkg/userinteraction"
)

func NewRunServerCmd() *cobra.Command {
	const short = "Runs the headscale as development docker container --container-name on --port"

	cmd := &cobra.Command{
		Use:   "run-server",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale localdev run-server`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.ContextVerbose()

			containerName, port := getContainerNameAndPort(cmd)

			_, cancel := mustutils.Must2(headscalelocaldevserver.RunLocalDevServer(ctx, &headscalelocaldevserver.RunOptions{
				Port:          port,
				ContainerName: containerName,
			}))
			defer cancel()

			userinteraction.WaitUserAbort(fmt.Sprintf("Container '%s' started. headscale development server is listening on port '%d'.", containerName, port))
		},
	}

	return cmd
}
