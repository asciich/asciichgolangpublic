package defaultclicommands

import (
	"github.com/spf13/cobra"
	dns_cmd "github.com/asciich/asciichgolangpublic/defaultclicommands/dns"
	errors_cmd "github.com/asciich/asciichgolangpublic/defaultclicommands/errors"
	gitlab_cmd "github.com/asciich/asciichgolangpublic/defaultclicommands/gitlab"
	monitoring_cmd "github.com/asciich/asciichgolangpublic/defaultclicommands/monitoring"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func AddDefaultCommands(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	const verbose_flag_name = "verbose"
	if rootCmd.PersistentFlags().Lookup(verbose_flag_name) == nil {
		rootCmd.PersistentFlags().Bool(verbose_flag_name, false, "Enable verbose output")
	}

	err = dns_cmd.AddDnsCommand(rootCmd)
	if err != nil {
		return err
	}

	err = errors_cmd.AddErrorsCommand(rootCmd)
	if err != nil {
		return err
	}

	err = gitlab_cmd.AddGitlabCommand(rootCmd)
	if err != nil {
		return err
	}

	err = monitoring_cmd.AddMonitoringGommand(rootCmd)
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
