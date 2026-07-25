package commandexecutorkubernetes

import (
	"context"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorDeployment struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func NewCommandExecutorDeployment() (c *CommandExecutorDeployment) {
	return new(CommandExecutorDeployment)
}

func (c *CommandExecutorDeployment) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorDeployment) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorDeployment) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorDeployment) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorDeployment) GetNamespaceName() (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorDeployment) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return nil, err
	}

	commandExecutorNamespace, ok := namespace.(*CommandExecutorNamespace)
	if !ok {
		typeName, _ := datatypes.GetTypeName(namespace)
		return nil, tracederrors.TracedErrorf("Only implemented for '*commandexecutorkubernetes.CommandExecutorNamespace' but got '%s'", typeName)
	}

	return commandExecutorNamespace.GetCommandExecutor()
}

func (c *CommandExecutorDeployment) GetKubectlContext(ctx context.Context) (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetKubectlContext(ctx)
}

func (c *CommandExecutorDeployment) Delete(ctx context.Context) error {
	deploymentName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete deployment '%s' in namespace '%s' started.", deploymentName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	deleteCommand := []string{
		"kubectl", "delete", "deployment", deploymentName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
	}

	var deleted bool
	var alreadyDeleted bool
	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: deleteCommand,
	})
	if err == nil {
		deleted = true
	} else {
		stderr, _ := output.GetStderrAsString()
		if strings.Contains(stderr, "Error from server (NotFound)") {
			deleted = false
			alreadyDeleted = true
		} else {
			return err
		}
	}

	if !alreadyDeleted {
		err := c.WaitUntilDeploymentDeleted(ctx)
		if err != nil {
			return err
		}
	}

	if deleted {
		logging.LogChangedByCtxf(ctx, "Deleted deployment '%s' in namespace '%s'.", deploymentName, namespaceName)
	} else {
		logging.LogChangedByCtxf(ctx, "Deployment '%s' in namespace '%s' is already absent. Skip delete.", deploymentName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Delete deployment '%s' in namespace '%s' finished.", deploymentName, namespaceName)

	return nil
}

func (c *CommandExecutorDeployment) WaitUntilDeploymentDeleted(ctx context.Context) error {
	deploymentName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until deployment '%s' in namespace '%s' is deleted started.", deploymentName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	getCommand := []string{
		"kubectl", "get", "deployment", deploymentName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
	}

	timeout := 5 * time.Minute
	interval := 5 * time.Second
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return tracederrors.TracedErrorf("timed out waiting for deployment '%s' in namespace '%s' to be deleted after %v", deploymentName, namespaceName, timeout)
		}

		output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: getCommand,
		})

		if err != nil {
			stderr, _ := output.GetStderrAsString()
			if strings.Contains(stderr, "Error from server (NotFound)") {
				break
			}
			return tracederrors.TracedErrorf("failed to check deployment '%s' in namespace '%s': %w", deploymentName, namespaceName, err)
		}

		logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' still exists. Waiting %v before retrying.", deploymentName, namespaceName, interval)

		select {
		case <-ctx.Done():
			return tracederrors.TracedErrorf("context cancelled while waiting for deployment '%s' in namespace '%s' to be deleted: %w", deploymentName, namespaceName, ctx.Err())
		case <-time.After(interval):
		}
	}

	logging.LogInfoByCtxf(ctx, "Wait until deployment '%s' in namespace '%s' is deleted finished.", deploymentName, namespaceName)

	return nil
}

func (c *CommandExecutorDeployment) Exists(ctx context.Context) (bool, error) {
	deploymentName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check if deployment '%s' in namespace '%s' exists started.", deploymentName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return false, err
	}

	getCommand := []string{
		"kubectl", "get", "deployment", deploymentName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
	}

	var exists bool
	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command:           getCommand,
		AllowAllExitCodes: true,
	})
	if err != nil {
		return false, err
	}

	if output.IsExitSuccess() {
		exists = true
	} else {
		stderr, err := output.GetStderrAsString()
		if err != nil {
			return false, err
		}

		if !strings.Contains(stderr, "Error from server (NotFound)") {
			return false, err
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' exists.", deploymentName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Deployment '%s' in namespace '%s' does not exist.", deploymentName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Check if deployment '%s' in namespace '%s' exists finished.", deploymentName, namespaceName)

	return exists, nil
}
