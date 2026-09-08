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

type CommandExecutorDaemonSet struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func NewCommandExecutorDaemonSet() (c *CommandExecutorDaemonSet) {
	return new(CommandExecutorDaemonSet)
}

func (c *CommandExecutorDaemonSet) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorDaemonSet) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorDaemonSet) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorDaemonSet) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorDaemonSet) GetNamespaceName() (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorDaemonSet) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
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

func (c *CommandExecutorDaemonSet) GetKubectlContext(ctx context.Context) (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetKubectlContext(ctx)
}

func (c *CommandExecutorDaemonSet) Delete(ctx context.Context) error {
	daemonSetName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete daemonset '%s' in namespace '%s' started.", daemonSetName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	deleteCommand := []string{
		"kubectl", "delete", "daemonset", daemonSetName,
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
		err := c.WaitUntilDaemonSetDeleted(ctx)
		if err != nil {
			return err
		}
	}

	if deleted {
		logging.LogChangedByCtxf(ctx, "Deleted daemonset '%s' in namespace '%s'.", daemonSetName, namespaceName)
	} else {
		logging.LogChangedByCtxf(ctx, "DaemonSet '%s' in namespace '%s' is already absent. Skip delete.", daemonSetName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Delete daemonset '%s' in namespace '%s' finished.", daemonSetName, namespaceName)

	return nil
}

func (c *CommandExecutorDaemonSet) WaitUntilDaemonSetDeleted(ctx context.Context) error {
	daemonSetName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until daemonset '%s' in namespace '%s' is deleted started.", daemonSetName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	getCommand := []string{
		"kubectl", "get", "daemonset", daemonSetName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
	}

	timeout := 5 * time.Minute
	interval := 5 * time.Second
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return tracederrors.TracedErrorf("timed out waiting for daemonset '%s' in namespace '%s' to be deleted after %v", daemonSetName, namespaceName, timeout)
		}

		output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: getCommand,
		})

		if err != nil {
			stderr, _ := output.GetStderrAsString()
			if strings.Contains(stderr, "Error from server (NotFound)") {
				break
			}
			return tracederrors.TracedErrorf("failed to check daemonset '%s' in namespace '%s': %w", daemonSetName, namespaceName, err)
		}

		logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' still exists. Waiting %v before retrying.", daemonSetName, namespaceName, interval)

		select {
		case <-ctx.Done():
			return tracederrors.TracedErrorf("context cancelled while waiting for daemonset '%s' in namespace '%s' to be deleted: %w", daemonSetName, namespaceName, ctx.Err())
		case <-time.After(interval):
		}
	}

	logging.LogInfoByCtxf(ctx, "Wait until daemonset '%s' in namespace '%s' is deleted finished.", daemonSetName, namespaceName)

	return nil
}

func (c *CommandExecutorDaemonSet) Exists(ctx context.Context) (bool, error) {
	daemonSetName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check if daemonset '%s' in namespace '%s' exists started.", daemonSetName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return false, err
	}

	getCommand := []string{
		"kubectl", "get", "daemonset", daemonSetName,
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
		logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' exists.", daemonSetName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "DaemonSet '%s' in namespace '%s' does not exist.", daemonSetName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Check if daemonset '%s' in namespace '%s' exists finished.", daemonSetName, namespaceName)

	return exists, nil
}
