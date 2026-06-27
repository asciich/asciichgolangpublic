package gen3handtcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/homeautomation/shelly/gen3handt"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunWebsocketServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-websocket-server",
		Short: "Run a websocket server on --port to receive the temperature and humidity for given --sensor-name.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			sensorNames, err := cmd.Flags().GetStringSlice("sensor-names")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if len(sensorNames) <= 0 {
				logging.LogFatal("Please specify --sensor-names")
			}

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if port <= 0 {
				logging.LogFatalf("Please specify a valid --port. Got '%d'.", port)
			}

			receiver := gen3handt.ShellyGen3HAndTWebsocketReceiver{
				Port:        port,
				SensorNames: sensorNames,
			}
			mustutils.Must0(receiver.Run(ctx))
		},
	}

	cmd.Flags().StringSlice("sensor-names", []string{}, "Name of the sensor.")
	cmd.Flags().Int("port", 0, "Port to run the wehook server")

	return cmd
}
