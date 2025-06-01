package defaultclicommands

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/ansiblecmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/dnscmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/errorscmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/gitlabcmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/loggingcmd"
	"github.com/asciich/asciichgolangpublic/defaultclicommands/monitoringcmd"
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

	rootCmd.AddCommand(
		ansiblecmd.NewAnsibleCmd(),
		dnscmd.NewDnsCommand(),
		errorscmd.NewErrorsCommand(),
		gitlabcmd.NewGitlabCommand(),
		loggingcmd.NewLoggingCmd(),
		monitoringcmd.NewMonitoringCommand(),
	)

	return nil
}
