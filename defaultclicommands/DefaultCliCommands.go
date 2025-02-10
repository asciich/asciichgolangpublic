package defaultclicommands

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/errors"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func AddDefaultCommands(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	err = errors.AddErrorsCommand(rootCmd)
	if err != nil {
		return err
	}

	return nil
}

func MustAddDefaultCommands(rootCmd *cobra.Command) {
	err := AddDefaultCommands(rootCmd)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}
