package errorscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/errorscmd/tracederrorscmd"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func NewErrorsCommand() (errorsCmd *cobra.Command) {
	errorsCmd = &cobra.Command{
		Use:   "errors",
		Short: "Error and Error handling related commands",
	}

	errorsCmd.AddCommand(
		tracederrorscmd.NewTracedErrorsCmd(),
	)

	return errorsCmd
}

func AddErrorsCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewErrorsCommand())

	return nil
}
