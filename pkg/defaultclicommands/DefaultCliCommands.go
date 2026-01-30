package defaultclicommands

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/ansiblecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/bashcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/dockercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/documentationcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/errorscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/filescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/installcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/latexcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/linuxcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/loggingcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/monitoringcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/packagemanagercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/sshcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/uuidcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/wikicmd"
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
		aicmd.NewAICmd(),
		ansiblecmd.NewAnsibleCmd(),
		bashcmd.NewBashCmd(),
		cloudcmd.NewCloudCmd(),
		dockercmd.NewDockerCmd(),
		documentationcmd.NewDocumentationCmd(rootCmd),
		errorscmd.NewErrorsCommand(),
		filescmd.NewFilesCmd(),
		gitlabcmd.NewGitlabCommand(),
		httpcmd.NewHttpCmd(),
		installcmd.NewInstallCmd(),
		kubernetescmd.NewKubernetesCmd(),
		latexcmd.NewLatexCmd(),
		linuxcmd.NewLinuxCmd(),
		loggingcmd.NewLoggingCmd(),
		monitoringcmd.NewMonitoringCommand(),
		networkcmd.NewNetworkCmd(),
		packagemanagercmd.NewPackageManagerCmd(),
		sshcmd.NewSshCmd(),
		uuidcmd.NewUuidCmd(),
		wikicmd.NewWikiCmd(),
	)

	return nil
}
