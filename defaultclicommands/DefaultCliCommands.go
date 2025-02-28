package defaultclicommands

import (
	"github.com/spf13/cobra"
	dns_cmd "github.com/asciich/asciichgolangpublic/defaultclicommands/dns"
	errors_cmd "github.com/asciich/asciichgolangpublic/defaultclicommands/errors"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func AddDefaultCommands(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")

	err = dns_cmd.AddDnsCommand(rootCmd)
	if err != nil {
		return err
	}

	err = errors_cmd.AddErrorsCommand(rootCmd)
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
