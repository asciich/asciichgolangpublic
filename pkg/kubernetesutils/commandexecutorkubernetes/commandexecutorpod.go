package commandexecutorkubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorPod struct {
	commandexecutorgeneric.CommandExecutorBase
	name      string
	namespace kubernetesinterfaces.Namespace
}

func NewCommandExecutorPod() *CommandExecutorPod {
	ret := new(CommandExecutorPod)
	ret.SetParentCommandExecutorForBaseClass(ret)
	return ret
}

func (c *CommandExecutorPod) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorPod) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorPod) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorPod) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {

	return c.namespace, nil
}

func (c *CommandExecutorPod) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorPod) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorPod) GetNamespaceName() (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorPod) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
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

func (c *CommandExecutorPod) GetKubectlContext(ctx context.Context) (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetKubectlContext(ctx)
}

func (c *CommandExecutorPod) Delete(ctx context.Context) error {
	podName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete pod '%s' in namespace '%s' started.", podName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	deleteCommand := []string{
		"kubectl", "delete", "pod", podName,
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
		err := c.WaitUntilPodDeleted(ctx)
		if err != nil {
			return err
		}
	}

	if deleted {
		logging.LogChangedByCtxf(ctx, "Deleted pod '%s' in namespace '%s'.", podName, namespaceName)
	} else {
		logging.LogChangedByCtxf(ctx, "Pod '%s' in namespace '%s' is already absent. Skip delete.", podName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Delete pod '%s' in namespace '%s' finished.", podName, namespaceName)

	return nil
}

func (c *CommandExecutorPod) WaitUntilPodDeleted(ctx context.Context) error {
	podName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until pod '%s' in namespace '%s' is deleted started.", podName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	getCommand := []string{
		"kubectl", "get", "pod", podName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
	}

	timeout := 5 * time.Minute
	interval := 5 * time.Second
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return tracederrors.TracedErrorf("timed out waiting for pod '%s' in namespace '%s' to be deleted after %v", podName, namespaceName, timeout)
		}

		output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: getCommand,
		})

		if err != nil {
			stderr, _ := output.GetStderrAsString()
			if strings.Contains(stderr, "Error from server (NotFound)") {
				break
			}
			return tracederrors.TracedErrorf("failed to check pod '%s' in namespace '%s': %w", podName, namespaceName, err)
		}

		logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' still exists. Waiting %v before retrying.", podName, namespaceName, interval)

		select {
		case <-ctx.Done():
			return tracederrors.TracedErrorf("context cancelled while waiting for pod '%s' in namespace '%s' to be deleted: %w", podName, namespaceName, ctx.Err())
		case <-time.After(interval):
		}
	}

	logging.LogInfoByCtxf(ctx, "Wait until pod '%s' in namespace '%s' is deleted finished.", podName, namespaceName)

	return nil
}

func (c *CommandExecutorPod) Exists(ctx context.Context) (bool, error) {
	podName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check if pod '%s' in namespace '%s' exists started.", podName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return false, err
	}

	getCommand := []string{
		"kubectl", "get", "pod", podName,
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
		logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' exists.", podName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Pod '%s' in namespace '%s' does not exist.", podName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Check if pod '%s' in namespace '%s' exists finished.", podName, namespaceName)

	return exists, nil
}

func (c *CommandExecutorPod) GetContainerLogs(ctx context.Context, containerName string) (stdout []byte, stderr []byte, err error) {
	podName, err := c.GetName()
	if err != nil {
		return nil, nil, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return nil, nil, err
	}

	logging.LogInfoByCtxf(ctx, "Get logs for container '%s' in pod '%s' in namespace '%s' started.", containerName, podName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, nil, err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get stdout logs
	stdoutCommand := []string{
		"kubectl", "logs", podName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
		"--container", containerName,
	}

	stdoutOutput, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: stdoutCommand,
	})
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("failed to get stdout logs for container '%s' in pod '%s' in namespace '%s': %w", containerName, podName, namespaceName, err)
	}

	stdoutBytes, err := stdoutOutput.GetStdoutAsBytes()
	if err != nil {
		return nil, nil, err
	}

	// Get stderr logs using --stderr flag (if supported) or return empty
	stderrCommand := []string{
		"kubectl", "logs", podName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
		"--container", containerName,
		"--stderr",
	}

	stderrOutput, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command:           stderrCommand,
		AllowAllExitCodes: true,
	})
	var stderrBytes []byte
	if err == nil && stderrOutput.IsExitSuccess() {
		stderrBytes, err = stderrOutput.GetStdoutAsBytes()
		if err != nil {
			return nil, nil, err
		}
	} else {
		// stderr flag not supported, return empty stderr
		stderrBytes = []byte{}
	}

	logging.LogInfoByCtxf(ctx, "Get logs for container '%s' in pod '%s' in namespace '%s' finished.", containerName, podName, namespaceName)

	return stdoutBytes, stderrBytes, nil
}

func (c *CommandExecutorPod) CopyFileToPod(ctx context.Context, localFile string, destPath string, containerName string) error {
	podName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Copy local file '%s' as '%s' into container '%s' of pod '%s' of namespace '%s' started.", localFile, destPath, containerName, podName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	// Use kubectl cp to copy the file
	copyCommand := []string{
		"kubectl", "cp", localFile,
		"--context", kubectlContext,
		"--namespace", namespaceName,
		podName + ":" + destPath,
		"-c", containerName,
	}

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: copyCommand,
	})
	if err != nil {
		return tracederrors.TracedErrorf("failed to copy file '%s' to pod '%s' in namespace '%s': %w", localFile, podName, namespaceName, err)
	}

	if !output.IsExitSuccess() {
		stderr, _ := output.GetStderrAsString()
		return tracederrors.TracedErrorf("kubectl cp failed for pod '%s' in namespace '%s': %s", podName, namespaceName, stderr)
	}

	logging.LogInfoByCtxf(ctx, "Copy local file '%s' as '%s' into container '%s' of pod '%s' of namespace '%s' finished.", localFile, destPath, containerName, podName, namespaceName)

	return nil
}

func (c *CommandExecutorPod) CopyFileFromPod(ctx context.Context, srcPath string, destFile string, containerName string) error {
	podName, err := c.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Copy file '%s' from container '%s' of pod '%s' of namespace '%s' to local '%s' started.", srcPath, containerName, podName, namespaceName, destFile)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return err
	}

	// Use kubectl cp to copy the file from pod
	copyCommand := []string{
		"kubectl", "cp",
		"--context", kubectlContext,
		"--namespace", namespaceName,
		podName + ":" + srcPath,
		destFile,
		"-c", containerName,
	}

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: copyCommand,
	})
	if err != nil {
		return tracederrors.TracedErrorf("failed to copy file from pod '%s' in namespace '%s': %w", podName, namespaceName, err)
	}

	if !output.IsExitSuccess() {
		stderr, _ := output.GetStderrAsString()
		return tracederrors.TracedErrorf("kubectl cp failed for pod '%s' in namespace '%s': %s", podName, namespaceName, stderr)
	}

	logging.LogInfoByCtxf(ctx, "Copy file '%s' from container '%s' of pod '%s' of namespace '%s' to local '%s' finished.", srcPath, containerName, podName, namespaceName, destFile)

	return nil
}

func (c *CommandExecutorPod) RunCommandInContainer(ctx context.Context, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (*commandoutput.CommandOutput, error) {
	podName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	containerName, err := options.GetContainerName()
	if err != nil {
		return nil, err
	}

	command, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command in pod '%s' container '%s' in namespace '%s' started.", podName, containerName, namespaceName)

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	// Use kubectl exec to run the command
	execCommand := []string{
		"kubectl", "exec", podName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
		"-c", containerName,
		"--",
	}
	execCommand = append(execCommand, command...)

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: execCommand,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to run command in pod '%s' in namespace '%s': %w", podName, namespaceName, err)
	}

	logging.LogInfoByCtxf(ctx, "Run command in pod '%s' container '%s' in namespace '%s' finished.", podName, containerName, namespaceName)

	return output, nil
}

func (c *CommandExecutorPod) GetDefaultContainerName(ctx context.Context) (string, error) {
	podName, err := c.GetName()
	if err != nil {
		return "", err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return "", err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return "", err
	}

	// Get pod details as JSON to extract container status information
	getCommand := []string{
		"kubectl", "get", "pod", podName,
		"--context", kubectlContext,
		"--namespace", namespaceName,
		"-o", "json",
	}

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: getCommand,
	})
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get pod '%s' in namespace '%s': %w", podName, namespaceName, err)
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get stdout from command: %w", err)
	}

	// Parse the JSON to find the first running container
	var podData map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &podData); err != nil {
		return "", tracederrors.TracedErrorf("Failed to parse pod JSON: %w", err)
	}

	// First try to get the container name from spec.containers (preferred for getting the defined name)
	spec, ok := podData["spec"].(map[string]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("Failed to get spec from pod JSON")
	}

	containers, ok := spec["containers"].([]interface{})
	if !ok || len(containers) == 0 {
		return "", tracederrors.TracedErrorf("Failed to get containers from pod spec")
	}

	// Get the first container name from spec
	// This gives us the actual container name as defined, not what kubectl might have set
	firstContainer, ok := containers[0].(map[string]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("Failed to parse first container in spec")
	}

	containerName, ok := firstContainer["name"].(string)
	if !ok || containerName == "" {
		return "", tracederrors.TracedErrorf("Failed to get container name from spec")
	}

	// Verify the container is actually running by checking status
	status, ok := podData["status"].(map[string]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("Failed to get status from pod JSON")
	}

	containerStatuses, ok := status["containerStatuses"].([]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("Failed to get containerStatuses from pod JSON")
	}

	// Check if any container is running
	hasRunningContainer := false
	for _, cs := range containerStatuses {
		containerStatus, ok := cs.(map[string]interface{})
		if !ok {
			continue
		}

		state, ok := containerStatus["state"].(map[string]interface{})
		if !ok {
			continue
		}

		// Check if the container is running
		if _, isRunning := state["running"]; isRunning {
			hasRunningContainer = true
			break
		}
	}

	if !hasRunningContainer {
		return "", tracederrors.TracedErrorf("No running container found in pod '%s' in namespace '%s'", podName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Default container name is '%s' in pod '%s'", containerName, podName)
	return containerName, nil
}

func (c *CommandExecutorPod) GetCPUArchitecture(ctx context.Context) (string, error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorPod) GetDeepCopyAsCommandExecutor() commandexecutorinterfaces.CommandExecutor {
	ret := NewCommandExecutorPod()

	*ret = *c

	return ret
}

func (c *CommandExecutorPod) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	podName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	containerName, err := c.GetDefaultContainerName(ctx)
	if err != nil {
		return nil, err
	}

	return c.RunCommandInContainer(ctx, &kubernetesparameteroptions.KubernetesRunCommandOptions{
		PodName:           podName,
		ContainerName:     containerName,
		RunCommandOptions: options,
	})
}

func (c *CommandExecutorPod) GetClusterName() (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetClusterName()
}

func (c *CommandExecutorPod) GetHostDescription() (string, error) {
	podName, err := c.GetName()
	if err != nil {
		return "", err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return "", err
	}

	clusterName, err := c.GetClusterName()
	if err != nil {
		return "", err
	}

	got := fmt.Sprintf("Pod '%s' in namespace '%s' of kubernetes cluster '%s'.", podName, namespaceName, clusterName)

	return got, nil
}

func (c *CommandExecutorPod) IsRunningOnLocalhost() (bool, error) {
	return false, nil
}
