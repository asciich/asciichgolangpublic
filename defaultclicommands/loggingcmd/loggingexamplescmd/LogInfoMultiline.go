package loggingexamplescmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/logging"
)

func NewLogInfoMultilineCmd() *cobra.Command {
	const logMessage = "This is the example log message\nwith multiple lines."
	const short = "Logs the multiline example message unsing 'logging.LogInfo'."
	long := fmt.Sprintf(
		"%s\nThe used example message is '%s'.",
		short,
		logMessage,
	)

	cmd := &cobra.Command{
		Use:   "log-info-multiline",
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			logging.LogInfo(logMessage)
		},
	}

	return cmd
}
