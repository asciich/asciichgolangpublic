package openwebuicmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/openwebuiutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunAsDockerContainerCmd() *cobra.Command {
	const short = "Runs openwebui as docker container on the local machine"

	cmd := &cobra.Command{
		Use:   "run-as-docker-container",
		Short: short,
		Long: short + `

Usage Example:
  ` + os.Args[0] + ` ai openwebui run-as-docker-container --container-name="openwebui" --port=8080

Example for a self hosted ollama on ollama.example.com:
  ` + os.Args[0] + ` ai openwebui run-as-docker-container --container-name="openwebui" --port=8080 --ollama-base-url="https://ollama.example.com:11434"
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

			dataPath, err := cmd.Flags().GetString("data-path")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			ollamaBaseUrl, err := cmd.Flags().GetString("ollama-base-url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			mustutils.Must(openwebuiutils.StartOpenWebUIAsDockerContainer(ctx, &openwebuiutils.StartOpenWebUIContainerOptions{
				Port:                     port,
				ContainerName:            containerName,
				ReachableByOtherMachines: reachableByOtherMachines,
				DataPath:                 dataPath,
				OllamaBaseUrl:            ollamaBaseUrl,
			}))

			logging.LogGoodByCtxf(ctx, "OpenWebUI container created")
		},
	}

	cmd.Flags().Int("port", 0, "Port for openwebui to listen to.")
	cmd.Flags().String("container-name", "", "Name of the docker container.")
	cmd.Flags().Bool("reachable-by-other-machines", false, "If set, binds to 0.0.0.0 instead of 127.0.0.1 to allow access from other machines in the network.")
	cmd.Flags().String("data-path", "", "Path to persist OpenWebUI data (optional).")
	cmd.Flags().String("ollama-base-url", "", "Base URL of Ollama instance to connect to (optional).")

	return cmd
}
