package errorscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/errorscmd/tracederrorscmd"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"os"
)

func NewErrorsCommand() (errorsCmd *cobra.Command) {
	const short = "Error and Error handling related commands"

	errorsCmd = &cobra.Command{
		Use:   "errors",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` errors errors`,
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
