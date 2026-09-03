package kubernetesutils_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// TestCopyFileToPod_nativeKubernetes demonstrates how to copy a local file to a container running in a Kubernetes pod
// using the nativekubernetes implementation. This is similar to the `kubectl cp` command.
func TestCopyFileToPod_nativeKubernetes(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	require.NoError(t, err)

	// Define the pod name and the namespace
	const podName = "example-copy-file-pod"
	const namespaceName = "example-copy-file"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Ensure the pod is absent before the example starts:
	_ = kubernetes.DeletePodByNames(ctx, namespaceName, podName)

	// Ensure pod is deleted at the end of the example:
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start the pod
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

	// Create a temporary file with some content
	testContent := "Hello from local file! This content will be copied to the container."
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	require.NoError(t, err)

	// Copy file to container
	destPath := "/tmp/example-file.txt"
	err = pod.CopyFileToPod(ctx, localPath, destPath, containerName)
	require.NoError(t, err)

	fmt.Printf("Successfully copied file to %s in pod %s\n", destPath, podName)

}

// ExampleCopyFileToPod_commandExecutor demonstrates how to copy a local file to a container
// using the commandexecutorkubernetes implementation (which uses kubectl cp under the hood).
func TestCopyFileToPod_commandExecutor(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := commandexecutorkubernetes.GetClusterByName("kind-" + kindutils.SharedClusterName)
	require.NoError(t, err)

	const podName = "example-copy-file-cmd-pod"
	const namespaceName = "example-copy-file-cmd"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Cleanup
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start the pod
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

	// Create a temporary file
	testContent := "File copied using kubectl cp implementation"
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	require.NoError(t, err)

	// Copy file to container using kubectl cp
	destPath := "/tmp/example-cmd-file.txt"
	err = pod.CopyFileToPod(ctx, localPath, destPath, containerName)
	require.NoError(t, err)

	fmt.Printf("Successfully copied file to %s in pod %s using kubectl cp\n", destPath, podName)

}

// ExampleCopyFileToPod_nestedDirectory demonstrates how to copy a file to a nested directory in a container.
func TestCopyFileToPod_nestedDirectory(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	require.NoError(t, err)

	const podName = "example-copy-nested-pod"
	const namespaceName = "example-copy-nested"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Cleanup
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start the pod
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
			WaitForPodRunning: true,
		},
	)
	require.NoError(t, err)

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	require.NoError(t, err)

	// Create a temporary file
	testContent := "File in nested directory"
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	require.NoError(t, err)

	// Copy file to a simple path (nested directories require pre-creation)
	destPath := "/tmp/example-nested-file.txt"
	err = pod.CopyFileToPod(ctx, localPath, destPath, containerName)
	require.NoError(t, err)

	fmt.Printf("Successfully copied file to %s in pod %s\n", destPath, podName)

}
