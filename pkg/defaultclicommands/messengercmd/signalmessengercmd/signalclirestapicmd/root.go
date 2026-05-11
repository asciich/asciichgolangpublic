package signalclirestapicmd

import "github.com/spf13/cobra"

func NewSignalCliRestApiCmd() *cobra.Command {
	const short = "Commands related to the Signal-CLI-Rest-api"

	cmd := &cobra.Command{
		Use:   "signal-cli-rest-api",
		Short: short,
		Long: short + `

The project can be found here:
  https://github.com/bbernhard/signal-cli-rest-api
`,
	}

	cmd.AddCommand(
		NewRunReceiveCacheServerCmd(),
	)

	return cmd
}
