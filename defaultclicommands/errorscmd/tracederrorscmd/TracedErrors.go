package tracederrorscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func NewTracedErrorsCommand() (errorsCmd *cobra.Command) {
	errorsCmd = &cobra.Command{
		Use:   "tracederrors",
		Short: "TracedErrors related commands",
	}

	errorsCmd.AddCommand(
		&cobra.Command{
			Use:   "demo",
			Short: "Demonstrate the output of a TracedError.",
			Run: func(cmd *cobra.Command, args []string) {
				cliDemo()
			},
		},
	)

	return errorsCmd
}

func AddTracedErrorsCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewTracedErrorsCommand())

	return nil
}

func cliDemo() {
	logging.LogGoErrorFatal(tracederrors.TracedError("Example TracedError"))
}
