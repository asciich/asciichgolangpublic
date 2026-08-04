package kubernetesutils_test

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

// ExampleCopyFileToPod demonstrates how to copy a local file to a container running in a Kubernetes pod
// using the nativekubernetes implementation. This is similar to the `kubectl cp` command.
func ExampleCopyFileToPod_nativeKubernetes() {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	// Define the pod name and the namespace
	const podName = "example-copy-file-pod"
	const namespaceName = "example-copy-file"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		panic(err)
	}

	// Ensure the pod is absent before the example starts:
	_ = kubernetes.DeletePodByNames(ctx, namespaceName, podName)

	// Ensure pod is deleted at the end of the example:
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start the pod
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			WaitForPodRunning:        true,
		},
	)
	if err != nil {
		panic(err)
	}

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	if err != nil {
		panic(err)
	}

	// Create a temporary file with some content
	testContent := "Hello from local file! This content will be copied to the container."
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	if err != nil {
		panic(err)
	}

	// Copy file to container
	destPath := "/tmp/example-file.txt"
	err = pod.CopyFileToPod(ctx, localPath, destPath, containerName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied file to %s in pod %s\n", destPath, podName)

	// Output: Successfully copied file to /tmp/example-file.txt in pod example-copy-file-pod
}

// ExampleCopyFileToPod_commandExecutor demonstrates how to copy a local file to a container
// using the commandexecutorkubernetes implementation (which uses kubectl cp under the hood).
func ExampleCopyFileToPod_commandExecutor() {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := commandexecutorkubernetes.GetClusterByName("kind-" + kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	const podName = "example-copy-file-cmd-pod"
	const namespaceName = "example-copy-file-cmd"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		panic(err)
	}

	// Cleanup
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start the pod
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			WaitForPodRunning:        true,
		},
	)
	if err != nil {
		panic(err)
	}

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	if err != nil {
		panic(err)
	}

	// Create a temporary file
	testContent := "File copied using kubectl cp implementation"
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	if err != nil {
		panic(err)
	}

	// Copy file to container using kubectl cp
	destPath := "/tmp/example-cmd-file.txt"
	err = pod.CopyFileToPod(ctx, localPath, destPath, containerName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied file to %s in pod %s using kubectl cp\n", destPath, podName)

	// Output: Successfully copied file to /tmp/example-cmd-file.txt in pod example-copy-file-cmd-pod using kubectl cp
}

// ExampleCopyFileToPod_nestedDirectory demonstrates how to copy a file to a nested directory in a container.
func ExampleCopyFileToPod_nestedDirectory() {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	const podName = "example-copy-nested-pod"
	const namespaceName = "example-copy-nested"
	const containerName = "test-container"

	// Create the namespace
	_, err = kubernetes.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		panic(err)
	}

	// Cleanup
	defer func() { _ = kubernetes.DeletePodByNames(ctx, namespaceName, podName) }()

	// Start the pod
	_, err = kubernetes.CreatePod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			ContainerName:            containerName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			WaitForPodRunning:        true,
		},
	)
	if err != nil {
		panic(err)
	}

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	if err != nil {
		panic(err)
	}

	// Create a temporary file
	testContent := "File in nested directory"
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	if err != nil {
		panic(err)
	}

	// Create the nested directory first using exec
	_, err = kubernetes.RunCommandInTemporaryPod(
		ctx,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  "mkdir-pod",
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"mkdir", "-p", "/tmp/nested/deep/path"},
		},
	)
	if err != nil {
		panic(err)
	}

	// Copy file to nested directory in container
	destPath := "/tmp/nested/deep/path/example-file.txt"
	err = pod.CopyFileToPod(ctx, localPath, destPath, containerName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied file to nested directory %s in pod %s\n", destPath, podName)

	// Output: Successfully copied file to nested directory /tmp/nested/deep/path/example-file.txt in pod example-copy-nested-pod
}
