package defaultclicommands

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/dnscmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/errorscmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/gitlabcmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/monitoringcmd"
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

	err = dnscmd.AddDnsCommand(rootCmd)
	if err != nil {
		return err
	}

	err = errorscmd.AddErrorsCommand(rootCmd)
	if err != nil {
		return err
	}

	err = gitlabcmd.AddGitlabCommand(rootCmd)
	if err != nil {
		return err
	}

	err = monitoringcmd.AddMonitoringGommand(rootCmd)
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
