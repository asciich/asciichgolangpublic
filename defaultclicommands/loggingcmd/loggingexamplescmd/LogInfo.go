package loggingexamplescmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
)

func NewLogInfoCmd() *cobra.Command {
	const logMessage = "This is the example log message."
	const short = "Logs the example message unsing 'logging.LogInfo'."
	long := fmt.Sprintf(
		"%s\nThe used example message is '%s'.",
		short,
		logMessage,
	)

	cmd := &cobra.Command{
		Use:   "log-info",
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			logging.LogInfo(logMessage)
		},
	}

	return cmd
}
