package errorscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/errorscmd/tracederrorscmd"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func NewErrorsCommand() (errorsCmd *cobra.Command) {
	errorsCmd = &cobra.Command{
		Use:   "errors",
		Short: "Error and Error handling related commands",
	}

	tracederrorscmd.AddTracedErrorsCommand(errorsCmd)

	return errorsCmd
}

func AddErrorsCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewErrorsCommand())

	return nil
}
