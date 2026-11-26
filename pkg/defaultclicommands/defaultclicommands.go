package defaultclicommands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/ansiblecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/certificatescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/cloudcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/containercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/datetimecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/documentationcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/embeddedsystemscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/errorscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/filescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/homeautomationcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/hostcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/installcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/kubernetescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/latexcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/linuxcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/loggingcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/messengercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/monitoringcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/packagemanagercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/shellcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/sshcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/templatescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/testingcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/uuidcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/versioncmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/virtualmachinescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/weatherutilscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/wikicmd"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitlabutils/gitlabinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitlabutils/gitlaboptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type genericGitlab struct{}

func (g *genericGitlab) GetGitlabInstance(ctx context.Context, cmd *cobra.Command, args []string) (gitlabinterfaces.GitlabInstance, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return nil, tracederrors.TracedError(err)
	}

	if url == "" {
		tracederrors.TracedError("Please specify the --url for the gitlab instance.")
	}

	return asciichgolangpublic.GetGitlabByFQDN(url)
}

func (g *genericGitlab) GetAuthenticatedGitlabInstance(ctx context.Context, cmd *cobra.Command, args []string) (gitlabinterfaces.GitlabInstance, error) {
	gitlabInstance, err := g.GetGitlabInstance(ctx, cmd, args)
	if err != nil {
		return nil, err
	}

	accessToken, err := environmentvariables.GetEnvValueAsString(ctx, "GITLAB_ACCESS_TOKEN", false)
	if err != nil {
		return nil, err
	}

	err = gitlabInstance.Authenticate(ctx, &gitlaboptions.GitlabAuthenticationOptions{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	return gitlabInstance, err
}

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
		certificatescmd.NewCertificatesCmd(),
		cloudcmd.NewCloudCmd(),
		containercmd.NewContainerCmd(),
		datetimecmd.NewDateTimeCmd(),
		documentationcmd.NewDocumentationCmd(rootCmd),
		embeddedsystemscmd.NewEmbeddedSystemsCmd(),
		errorscmd.NewErrorsCommand(),
		filescmd.NewFilesCmd(),
		gitlabcmd.NewGitlabCommand(&genericGitlab{}),
		homeautomationcmd.NewHomeAutomationCmd(),
		hostcmd.NewHostCmd(),
		httpcmd.NewHttpCmd(),
		installcmd.NewInstallCmd(),
		kubernetescmd.NewKubernetesCmd(),
		latexcmd.NewLatexCmd(),
		linuxcmd.NewLinuxCmd(),
		loggingcmd.NewLoggingCmd(),
		messengercmd.NewMessengerCmd(),
		monitoringcmd.NewMonitoringCommand(),
		networkcmd.NewNetworkCmd(),
		packagemanagercmd.NewPackageManagerCmd(),
		shellcmd.NewShellCmd(),
		sshcmd.NewSshCmd(),
		storagecmd.NewStorageCmd(),
		templatescmd.NewTemplatesCmd(),
		testingcmd.NewTestingCmd(),
		uuidcmd.NewUuidCmd(),
		versioncmd.NewVersionCmd(),
		virtualmachinescmd.NewVirtualMachinesCmd(),
		weatherutilscmd.NewWeatherUtilsCmd(),
		wikicmd.NewWikiCmd(),
	)

	return nil
}
