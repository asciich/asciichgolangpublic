package errors

import (
	"github.com/spf13/cobra"
	tracederrors_cmd "github.com/asciich/asciichgolangpublic/defaultclicommands/errors/tracederrors"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func NewErrorsCommand() (errorsCmd *cobra.Command) {
	errorsCmd = &cobra.Command{
		Use:   "errors",
		Short: "Error and Error handling related commands",
	}

	tracederrors_cmd.AddTracedErrorsCommand(errorsCmd)

	return errorsCmd
}

func AddErrorsCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewErrorsCommand())

	return nil
}
