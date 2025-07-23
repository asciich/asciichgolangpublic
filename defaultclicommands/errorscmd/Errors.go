package errorscmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/defaultclicommands/errorscmd/tracederrorscmd"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
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
