package openhandscmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/openhandsutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunAsDockerContainerCmd() *cobra.Command {
	const short = "Runs openhands as docker container on the local machine"

	cmd := &cobra.Command{
		Use:   "run-as-docker-container",
		Short: short,
		Long: short + `

Usage:
  ` + os.Args[0] + `ai openhands run-as-docker-container --port=8000 --container-name=openhands --reachable-by-other-machines --verbose
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if port <= 0 {
				logging.LogFatalf("Invalid port '%d'.", port)
			}

			containerName, err := cmd.Flags().GetString("container-name")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if containerName == "" {
				logging.LogFatal("Please specify --container-name")
			}

			reachableByOtherMachines, err := cmd.Flags().GetBool("reachable-by-other-machines")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			workspacePath, err := cmd.Flags().GetString("workspace-path")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if workspacePath == "" {
				logging.LogFatal("Please specify --workspace-path")
			}

			mustutils.Must(openhandsutils.StartAsDockerContainer(ctx, &openhandsutils.StartContainerOptions{
				Port:                     port,
				ContainerName:            containerName,
				ReachableByOtherMachines: reachableByOtherMachines,
				WorkspacePath:            workspacePath,
			}))

			logging.LogGoodByCtxf(ctx, "Openhands container created")
		},
	}

	cmd.Flags().Int("port", 0, "Port for openhands to listen to.")
	cmd.Flags().String("container-name", "", "Name of the docker container.")
	cmd.Flags().Bool("reachable-by-other-machines", false, "If set, binds to 0.0.0.0 instead of 127.0.0.1 to allow access from other machines in the network.")
	cmd.Flags().String("workspace-path", ".", "Path to the workspace directory to mount into the container.")

	return cmd
}
