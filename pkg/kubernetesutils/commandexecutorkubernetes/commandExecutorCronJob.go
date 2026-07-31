package commandexecutorkubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorCronJob struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func (c *CommandExecutorCronJob) GetName() (string, error) {
	if c.name == "" {
		return "", tracederrors.TracedError("name not set")
	}
	return c.name, nil
}

func (c *CommandExecutorCronJob) GetCommandExecutorNamespace() (*CommandExecutorNamespace, error) {
	if c.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}
	namespace, ok := c.namespace.(*CommandExecutorNamespace)
	if !ok {
		return nil, tracederrors.TracedErrorf("namespace is not of type *CommandExecutorNamespace but '%T'", c.namespace)
	}
	return namespace, nil
}

func (c *CommandExecutorCronJob) GetNamespace() (kubernetesinterfaces.Namespace, error) {
	if c.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}
	return c.namespace, nil
}

func (c *CommandExecutorCronJob) Exists(ctx context.Context) (bool, error) {
	cronJobName, err := c.GetName()
	if err != nil {
		return false, err
	}
	namespace, err := c.GetNamespace()
	if err != nil {
		return false, err
	}
	return namespace.CronJobByNameExists(ctx, cronJobName)
}

func (c *CommandExecutorCronJob) GetKubernetesCluster() (kubernetesinterfaces.KubernetesCluster, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return nil, err
	}
	return namespace.GetKubernetesCluster()
}

func (c *CommandExecutorCronJob) GetCachedKubectlContext(ctx context.Context) (string, error) {
	namespace, err := c.GetCommandExecutorNamespace()
	if err != nil {
		return "", err
	}
	return namespace.GetCachedKubectlContext(ctx)
}

func (c *CommandExecutorCronJob) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	namespace, err := c.GetCommandExecutorNamespace()
	if err != nil {
		return nil, err
	}
	return namespace.RunCommand(ctx, options)
}

func (c *CommandExecutorCronJob) RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]string, error) {
	namespace, err := c.GetCommandExecutorNamespace()
	if err != nil {
		return nil, err
	}
	return namespace.RunCommandAndGetStdoutAsLines(ctx, options)
}
