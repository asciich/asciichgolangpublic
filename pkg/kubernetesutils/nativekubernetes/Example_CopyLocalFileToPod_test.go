package nativekubernetes_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// ExampleCopyFileToPod demonstrates how to copy a local file to a container running in a Kubernetes pod.
// This is similar to the `kubectl cp` command.
func Test_ExampleCopyFileToPod(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Define the pod name and the namespace
	const podName = "example-copy-local-file-pod"
	namespaceName := "example-copy-local-file"

	// Get the cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	if err != nil {
		panic(err)
	}

	// Create the namespace
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		panic(err)
	}

	// Get the config and clientset to access the kubernetes cluster:
	config, err := nativekubernetes.GetConfig(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	// Ensure the pod is absent before the example starts:
	err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	if err != nil {
		panic(err)
	}

	// Ensure pod is deleted at the end of the example:
	defer func() { _ = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName) }()

	// Start the pod
	err = nativekubernetes.CreatePod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
			WaitForPodRunning: true,
		},
	)
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
	containerName := podName
	err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied file to %s in pod %s\n", destPath, podName)

	// Verify the file was copied correctly by reading it back
	commandOutput, err := nativekubernetes.Exec(
		ctx,
		config,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			PodName:       podName,
			ContainerName: containerName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"cat", destPath},
			},
		},
	)
	if err != nil {
		panic(err)
	}

	readContent, err := commandOutput.GetStdoutAsString()
	if err != nil {
		panic(err)
	}

	if readContent == testContent {
		fmt.Println("File content verified successfully!")
	}
}

// ExampleCopyFileToPod_binary demonstrates how to copy a binary file to a container.
func Test_ExampleCopyFileToPod_binary(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	const podName = "example-copy-binary-pod"
	namespaceName := "example-copy-binary"

	// Get the cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	if err != nil {
		panic(err)
	}

	// Create the namespace
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		panic(err)
	}

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	// Cleanup
	defer func() {
		_ = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	}()

	// Start the pod
	err = nativekubernetes.CreatePod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
			WaitForPodRunning: true,
		},
	)
	if err != nil {
		panic(err)
	}

	// Create a temporary binary file
	binaryContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG magic number
	localPath, err := tempfiles.CreateTemporaryFileFromContentBytes(ctx, binaryContent)
	if err != nil {
		panic(err)
	}

	// Copy binary file to container
	destPath := "/tmp/example-binary.bin"
	containerName := podName
	err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied binary file to %s in pod %s\n", destPath, podName)
}

// ExampleCopyFileToPod_nestedDirectory demonstrates how to copy a file to a nested directory in a container.
func Test_ExampleCopyFileToPod_nestedDirectory(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	const podName = "example-copy-nested-pod"
	namespaceName := "example-copy-nested"

	// Get the cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	if err != nil {
		panic(err)
	}

	// Create the namespace
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		panic(err)
	}

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+kindutils.SharedClusterName)
	if err != nil {
		panic(err)
	}

	// Cleanup
	defer func() {
		_ = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	}()

	// Start the pod
	err = nativekubernetes.CreatePod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			},
			WaitForPodRunning: true,
		},
	)
	if err != nil {
		panic(err)
	}

	// Create a temporary file
	testContent := "File in nested directory"
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	if err != nil {
		panic(err)
	}

	// Create the nested directory first
	_, err = nativekubernetes.Exec(
		ctx,
		config,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			PodName: podName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"mkdir", "-p", "/tmp/nested/deep/path"},
			},
			ContainerName: podName,
		},
	)
	if err != nil {
		panic(err)
	}

	// Copy file to nested directory in container
	destPath := "/tmp/nested/deep/path/example-file.txt"
	containerName := podName
	err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully copied file to nested directory %s in pod %s\n", destPath, podName)
}
