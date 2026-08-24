package gen3handtcmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewGen3HAndTCmd() *cobra.Command {
	const short = "Gen3 H&T Humidity and Temperature sensor."

	cmd := &cobra.Command{
		Use:   "gen3-h-and-t",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` homeautomation shelly gen3handt gen3-h-and-t`,
	}

	cmd.AddCommand(
		NewRunWebsocketServerCmd(),
	)

	return cmd
}
