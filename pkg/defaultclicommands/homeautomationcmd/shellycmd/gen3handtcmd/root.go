package gen3handtcmd

import "github.com/spf13/cobra"

func NewGen3HAndTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen3-h-and-t",
		Short: "Gen3 H&T Humidity and Temperature sensor.",
	}

	cmd.AddCommand(
		NewRunWebsocketServerCmd(),
	)

	return cmd
}
