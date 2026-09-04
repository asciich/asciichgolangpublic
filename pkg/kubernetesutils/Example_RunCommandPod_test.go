package kubernetesutils_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// ExamplePod_RunCommand_nativeKubernetes demonstrates how to run a command in an existing pod
// using the nativekubernetes implementation (using the Kubernetes API directly).
func TestPod_RunCommand_nativeKubernetes(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	require.NoError(t, err)

	// Define the pod name and the namespace
	const podName = "example-run-command-pod"
	const namespaceName = "example-run-command"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Ensure the pod is absent before the example starts:
	_ = kubernetes.DeletePodByNames(ctx, namespaceName, podName)

	// Ensure pod is deleted at the end of the example:
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start a long-running pod that we can execute commands in
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			WaitForPodRunning:        true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
		},
	)
	require.NoError(t, err)

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	require.NoError(t, err)

	// Run a command in the existing pod
	output, err := pod.RunCommandInContainer(
		ctx,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			ContainerName: containerName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "echo 'Hello from existing pod!'"},
			},
		},
	)
	require.NoError(t, err)

	// Read the stdout of the executed command
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)

	fmt.Printf("Command output: %s", stdout)

}

// ExamplePod_RunCommand_commandExecutor demonstrates how to run a command in an existing pod
// using the commandexecutorkubernetes implementation (which uses kubectl exec under the hood).
func TestPod_RunCommand_commandExecutor(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := commandexecutorkubernetes.GetClusterByName("kind-" + kindutils.SharedClusterName)
	require.NoError(t, err)

	const podName = "example-run-command-cmd-pod"
	const namespaceName = "example-run-command-cmd"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Cleanup
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start a long-running pod
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			WaitForPodRunning:        true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
		},
	)
	require.NoError(t, err)

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	require.NoError(t, err)

	// Run a command in the existing pod using kubectl exec
	output, err := pod.RunCommandInContainer(
		ctx,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			ContainerName: containerName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "echo 'Hello from pod via kubectl exec!'"},
			},
		},
	)
	require.NoError(t, err)

	// Read the stdout
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)

	fmt.Printf("Command output via kubectl exec: %s", stdout)

}

// ExamplePod_RunCommand_multipleCommands demonstrates how to run multiple commands in the same pod.
func TestPod_RunCommand_multipleCommands(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	require.NoError(t, err)

	const podName = "example-multi-cmd-pod"
	const namespaceName = "example-multi-cmd"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Cleanup
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start a long-running pod
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			WaitForPodRunning:        true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
		},
	)
	require.NoError(t, err)

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	require.NoError(t, err)

	// Run first command: get hostname
	output1, err := pod.RunCommandInContainer(
		ctx,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			ContainerName: containerName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"hostname"},
			},
		},
	)
	require.NoError(t, err)

	stdout1, err := output1.GetStdoutAsString()
	require.NoError(t, err)

	// Run second command: echo a fixed string
	output2, err := pod.RunCommandInContainer(
		ctx,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			ContainerName: containerName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"echo", "multiple_commands_work"},
			},
		},
	)
	require.NoError(t, err)

	stdout2, err := output2.GetStdoutAsString()
	require.NoError(t, err)

	fmt.Printf("Hostname: %s", stdout1)
	fmt.Printf("Second command: %s", stdout2)

	// Hostname: example-multi-cmd-pod
	// Second command: multiple_commands_work
}

// ExamplePod_RunCommand_withExitCode demonstrates how to run a command and check its exit code.
func TestPod_RunCommand_withExitCode(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	require.NoError(t, err)

	const podName = "example-exit-code-pod"
	const namespaceName = "example-exit-code"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Cleanup
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start a long-running pod
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			WaitForPodRunning:        true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
		},
	)
	require.NoError(t, err)

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	require.NoError(t, err)

	// Run a command that succeeds
	output1, err := pod.RunCommandInContainer(
		ctx,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			ContainerName: containerName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "echo success && exit 0"},
			},
		},
	)
	require.NoError(t, err)

	exitCode1, err := output1.GetReturnCode()
	require.NoError(t, err)

	// Run a command that fails
	output2, err := pod.RunCommandInContainer(
		ctx,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			ContainerName: containerName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command:           []string{"bash", "-c", "echo error && exit 1"},
				AllowAllExitCodes: true, // Otherwise an error would be returned.
			},
		},
	)
	require.NoError(t, err)

	exitCode2, err := output2.GetReturnCode()

	fmt.Printf("First command exit code: %d\n", exitCode1)
	fmt.Printf("Second command exit code: %d\n", exitCode2)

	// First command exit code: 0
	// Second command exit code: 1
}
