package helmutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils/helminterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils/helmparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type commandExecutorHelm struct {
	commandExecutor commandexecutor.CommandExecutor
}

func GetCommandExecutorHelm(executor commandexecutor.CommandExecutor) (helm helminterfaces.Helm, err error) {
	if executor == nil {
		return nil, tracederrors.TracedErrorNil("executor")
	}

	toReturn := NewcommandExecutorHelm()

	err = toReturn.SetCommandExecutor(executor)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetLocalCommandExecutorHelm() (helm helminterfaces.Helm, err error) {
	return GetCommandExecutorHelm(commandexecutor.Bash())
}

func NewcommandExecutorHelm() (c *commandExecutorHelm) {
	return new(commandExecutorHelm)
}

func (c *commandExecutorHelm) AddRepositoryByName(ctx context.Context, name string, url string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if url == "" {
		return tracederrors.TracedErrorEmptyString("url")
	}

	commandExecutor, hostDescription, err := c.GetCommandExecutorAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add helm repository '%s' with url '%s' on host '%s' started.", name, url, hostDescription)

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"helm",
				"repo",
				"add",
				name,
				url,
			},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Added helm repository '%s' with url '%s' on host '%s'.", name, url, hostDescription)
	logging.LogInfoByCtxf(ctx, "Add helm repository '%s' with url '%s' on host '%s' finished.", name, url, hostDescription)

	return nil
}

func (c *commandExecutorHelm) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *commandExecutorHelm) GetCommandExecutorAndHostDescription() (commandExecutor commandexecutor.CommandExecutor, hostDescription string, err error) {
	commandExecutor, err = c.GetCommandExecutor()
	if err != nil {
		return nil, "", err
	}

	hostDescription, err = c.GetHostDescription()
	if err != nil {
		return nil, "", err
	}

	return commandExecutor, hostDescription, nil
}

func (c *commandExecutorHelm) GetHostDescription() (hostDescription string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.GetHostDescription()
}

func (c *commandExecutorHelm) SetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *commandExecutorHelm) InstallHelmChart(ctx context.Context, options *helmparameteroptions.InstallHelmChartOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	cluster, err := options.GetKubernetesCluster()
	if err != nil {
		return err
	}

	kubeContext, err := cluster.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	chartReference, err := options.GetChartReference()
	if err != nil {
		return err
	}

	chartUri, err := options.GetChartUri()
	if err != nil {
		return err
	}

	namespace, err := options.GetNamespace()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install helm chart '%s' as '%s' in namespace '%s' using kube context '%s' started.", chartUri, chartReference, namespace, kubeContext)

	cmd := []string{"helm", "upgrade", "--install", "--kube-context", kubeContext, chartReference, chartUri, "--namespace", namespace, "--create-namespace", "--wait"}
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		commandexecutor.WithLiveOutputOnStdout(ctx),
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install helm chart '%s' as '%s' in namespace '%s' using kube context '%s' finished.", chartUri, chartReference, namespace, kubeContext)

	return nil
}
