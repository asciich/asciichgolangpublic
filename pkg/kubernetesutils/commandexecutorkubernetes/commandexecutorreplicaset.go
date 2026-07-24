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

type CommandExecutorReplicaSet struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func NewCommandExecutorReplicaSet() (c *CommandExecutorReplicaSet) {
	return new(CommandExecutorReplicaSet)
}

func (c *CommandExecutorReplicaSet) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorReplicaSet) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorReplicaSet) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorReplicaSet) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorReplicaSet) GetNamespaceName() (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorReplicaSet) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
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

func (c *CommandExecutorReplicaSet) GetKubectlContext(ctx context.Context) (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetKubectlContext(ctx)
}

func (c *CommandExecutorReplicaSet) Delete(ctx context.Context) error {
	replicaSetName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete replicaset '%s' in namespace '%s' started.", replicaSetName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	deleteCommand := []string{
		"kubectl", "delete", "replicaset", replicaSetName,
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
		err := c.WaitUntilReplicaSetDeleted(ctx)
		if err != nil {
			return err
		}
	}

	if deleted {
		logging.LogChangedByCtxf(ctx, "Deleted replicaset '%s' in namespace '%s'.", replicaSetName, namespaceName)
	} else {
		logging.LogChangedByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' is already absent. Skip delete.", replicaSetName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Delete replicaset '%s' in namespace '%s' finished.", replicaSetName, namespaceName)

	return nil
}

func (c *CommandExecutorReplicaSet) WaitUntilReplicaSetDeleted(ctx context.Context) error {
	replicaSetName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until replicaset '%s' in namespace '%s' is deleted started.", replicaSetName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	getCommand := []string{
		"kubectl", "get", "replicaset", replicaSetName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
	}

	timeout := 5 * time.Minute
	interval := 5 * time.Second
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return tracederrors.TracedErrorf("timed out waiting for replicaset '%s' in namespace '%s' to be deleted after %v", replicaSetName, namespaceName, timeout)
		}

		output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: getCommand,
		})

		if err != nil {
			stderr, _ := output.GetStderrAsString()
			if strings.Contains(stderr, "Error from server (NotFound)") {
				break
			}
			return tracederrors.TracedErrorf("failed to check replicaset '%s' in namespace '%s': %w", replicaSetName, namespaceName, err)
		}

		logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' still exists. Waiting %v before retrying.", replicaSetName, namespaceName, interval)

		select {
		case <-ctx.Done():
			return tracederrors.TracedErrorf("context cancelled while waiting for replicaset '%s' in namespace '%s' to be deleted: %w", replicaSetName, namespaceName, ctx.Err())
		case <-time.After(interval):
		}
	}

	logging.LogInfoByCtxf(ctx, "Wait until replicaset '%s' in namespace '%s' is deleted finished.", replicaSetName, namespaceName)

	return nil
}

func (c *CommandExecutorReplicaSet) Exists(ctx context.Context) (bool, error) {
	replicaSetName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check if replicaset '%s' in namespace '%s' exists started.", replicaSetName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return false, err
	}

	getCommand := []string{
		"kubectl", "get", "replicaset", replicaSetName,
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
		logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' exists.", replicaSetName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "ReplicaSet '%s' in namespace '%s' does not exist.", replicaSetName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Check if replicaset '%s' in namespace '%s' exists finished.", replicaSetName, namespaceName)

	return exists, nil
}
