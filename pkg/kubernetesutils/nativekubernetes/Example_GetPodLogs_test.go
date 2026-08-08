package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// This example shows how to fetch logs from a container running in a k8s pod.
func Test_Example_GetPodLogs(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// define the pod name and the namespace
	const podName = "example-logs-pod"
	const containerName = "example-container"
	namespaceName := "test-example-getpodlogs"

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	// Get the config and clientset to access the kubernetes cluster:
	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// Ensure the pod is absent before the test starts:
	err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)

	// Ensure pod is deleted at the end of the test:
	defer func() { nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName) }()

	// Check pod absent
	exists, err := nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// start the pod that writes to stdout and stderr
	err = nativekubernetes.CreatePod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			PodName:                  podName,
			ContainerName:            containerName,
			Image:                    "ubuntu",
			DeleteAlreadyExistingPod: true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "echo 'Hello from stdout'; echo 'Error from stderr' >&2; sleep 10"},
			},
			WaitForPodRunning: true,
		},
	)
	require.NoError(t, err)

	// Check pod present
	exists, err = nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Get stdout logs (may not be supported in all Kubernetes APIs)
	stdout, err := nativekubernetes.GetContainerStdoutLogs(ctx, clientset, namespaceName, podName, containerName)
	if err != nil {
		// Some Kubernetes APIs don't support separate stream specification
		// Fall back to combined logs
		t.Logf("Stdout logs not separately available: %v, using combined logs", err)
		allLogs, err := nativekubernetes.GetContainerAllLogs(ctx, clientset, namespaceName, podName, containerName)
		require.NoError(t, err)
		require.Contains(t, string(allLogs), "Hello from stdout")
	} else {
		require.Contains(t, string(stdout), "Hello from stdout")
	}

	// Get stderr logs (may not be supported in all Kubernetes APIs)
	stderr, err := nativekubernetes.GetContainerStderrLogs(ctx, clientset, namespaceName, podName, containerName)
	if err != nil {
		// Some Kubernetes APIs don't support separate stream specification
		t.Logf("Stderr logs not separately available: %v", err)
	} else {
		require.Contains(t, string(stderr), "Error from stderr")
	}

	// Get combined logs (stdout + stderr) - always available
	allLogs, err := nativekubernetes.GetContainerAllLogs(ctx, clientset, namespaceName, podName, containerName)
	require.NoError(t, err)
	require.Contains(t, string(allLogs), "Hello from stdout")
	require.Contains(t, string(allLogs), "Error from stderr")

	// Get both stdout and stderr separately using GetContainerLogs
	// This function handles the fallback automatically when separate streams aren't available
	stdout2, _, err := nativekubernetes.GetContainerLogs(ctx, clientset, namespaceName, podName, containerName)
	require.NoError(t, err)
	require.Contains(t, string(stdout2), "Hello from stdout")
}
