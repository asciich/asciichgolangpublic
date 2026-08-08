package kubernetesutils_test

import (
	"context"
	"fmt"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// ExampleCopyFileFromPod demonstrates how to copy a file from a container running in a Kubernetes pod
// to the local filesystem using the nativekubernetes implementation. This is similar to the `kubectl cp` command.
func ExampleCopyFileFromPod_nativeKubernetes() {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	// Define the pod name and the namespace
	const podName = "example-copy-from-pod"
	const namespaceName = "example-copy-from"
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
	if err != nil {
		panic(err)
	}

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	if err != nil {
		panic(err)
	}

	// First, create a file in the pod by copying TO it
	testContent := "Hello from local file! This content will be copied to the container and back."
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	if err != nil {
		panic(err)
	}

	// Copy file to container first
	srcPathInPod := "/tmp/example-source-file.txt"
	err = pod.CopyFileToPod(ctx, localPath, srcPathInPod, containerName)
	if err != nil {
		panic(err)
	}

	// Now copy the file back from the container to local filesystem
	destFile, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		panic(err)
	}

	err = pod.CopyFileFromPod(ctx, srcPathInPod, destFile, containerName)
	if err != nil {
		panic(err)
	}

	// Read and verify the copied content
	content, err := os.ReadFile(destFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied file from %s in pod %s to local filesystem. Content matches: %v\n", srcPathInPod, podName, string(content) == testContent)

	// Output: Successfully copied file from /tmp/example-source-file.txt in pod example-copy-from-pod to local filesystem. Content matches: true
}

// ExampleCopyFileFromPod_commandExecutor demonstrates how to copy a file from a container
// using the commandexecutorkubernetes implementation (which uses kubectl cp under the hood).
func ExampleCopyFileFromPod_commandExecutor() {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := commandexecutorkubernetes.GetClusterByName("kind-" + kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	const podName = "example-copy-from-cmd-pod"
	const namespaceName = "example-copy-from-cmd"
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
	if err != nil {
		panic(err)
	}

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	if err != nil {
		panic(err)
	}

	// First, create a file in the pod
	testContent := "File copied using kubectl cp implementation"
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	if err != nil {
		panic(err)
	}

	// Copy file to pod first
	srcPathInPod := "/tmp/example-cmd-source.txt"
	err = pod.CopyFileToPod(ctx, localPath, srcPathInPod, containerName)
	if err != nil {
		panic(err)
	}

	// Now copy it back from the pod using kubectl cp
	destFile, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		panic(err)
	}

	err = pod.CopyFileFromPod(ctx, srcPathInPod, destFile, containerName)
	if err != nil {
		panic(err)
	}

	// Read and verify the copied content
	content, err := os.ReadFile(destFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied file from %s in pod %s to local filesystem using kubectl cp. Content matches: %v\n", srcPathInPod, podName, string(content) == testContent)

	// Output: Successfully copied file from /tmp/example-cmd-source.txt in pod example-copy-from-cmd-pod to local filesystem using kubectl cp. Content matches: true
}

// ExampleCopyFileFromPod_roundTrip demonstrates a complete round-trip: copy a file to a pod and back.
func ExampleCopyFileFromPod_roundTrip() {
	ctx := contextutils.WithVerbose(context.TODO())

	// Get the Kubernetes implementation
	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	const podName = "example-roundtrip-pod"
	const namespaceName = "example-roundtrip"
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
	if err != nil {
		panic(err)
	}

	// Get the pod object
	pod, err := kubernetes.GetPodByNames(namespaceName, podName)
	if err != nil {
		panic(err)
	}

	// Create original file with test content
	originalContent := "Round-trip test content\nLine 2\nLine 3 with special chars: äöü"
	originalFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, originalContent)
	if err != nil {
		panic(err)
	}

	// Step 1: Copy file TO the pod
	podPath := "/tmp/roundtrip-file.txt"
	err = pod.CopyFileToPod(ctx, originalFile, podPath, containerName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Step 1: Copied file to pod at %s\n", podPath)

	// Step 2: Copy file FROM the pod back to local
	retrievedFile, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		panic(err)
	}

	err = pod.CopyFileFromPod(ctx, podPath, retrievedFile, containerName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Step 2: Copied file from pod to local filesystem\n")

	// Step 3: Verify content integrity
	retrievedContent, err := os.ReadFile(retrievedFile)
	if err != nil {
		panic(err)
	}

	contentMatches := string(retrievedContent) == originalContent
	fmt.Printf("Step 3: Content integrity check: %v\n", contentMatches)

	// Output:
	// Step 1: Copied file to pod at /tmp/roundtrip-file.txt
	// Step 2: Copied file from pod to local filesystem
	// Step 3: Content integrity check: true
}
