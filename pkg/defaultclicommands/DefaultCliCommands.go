package defaultclicommands

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/ansiblecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/dnscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/dockercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/errorscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/filescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/loggingcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/monitoringcmd"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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
		dockercmd.NewDockerCmd(),
		errorscmd.NewErrorsCommand(),
		filescmd.NewFilesCmd(),
		gitlabcmd.NewGitlabCommand(),
		kubernetescmd.NewKubernetesCmd(),
		loggingcmd.NewLoggingCmd(),
		monitoringcmd.NewMonitoringCommand(),
	)

	return nil
}
