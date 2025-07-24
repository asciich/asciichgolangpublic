package loggingexamplescmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func NewLogInfoMultilineCmd() *cobra.Command {
	const linePrefix = "lineprefix"

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
			lineprefix, err := cmd.Flags().GetBool("lineprefix")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if lineprefix {
				ctx := contextutils.ContextVerbose()
				ctx = contextutils.WithLogLinePrefix(ctx, linePrefix)
				logging.LogInfoByCtx(ctx, logMessage)
			} else {
				logging.LogInfo(logMessage)
			}
		},
	}

	cmd.PersistentFlags().Bool(
		"lineprefix",
		false,
		fmt.Sprintf("Additionally use line prefix '%s'.", linePrefix),
	)

	return cmd
}
