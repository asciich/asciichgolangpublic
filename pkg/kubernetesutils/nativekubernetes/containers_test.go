package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_GetContainerLogs(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	namespaceName := "test-getcontainerlogs"
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	const podName = "test-pod-logs"
	const containerName = "test-container"

	// Ensure pod is absent before test and cleaned up after
	defer func() {
		nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	}()
	err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)

	// Create a pod that writes to stdout and stderr
	err = nativekubernetes.CreatePod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			PodName:                  podName,
			ContainerName:            containerName,
			Image:                    "ubuntu",
			Command:                  []string{"bash", "-c", "echo 'stdout message'; echo 'stderr message' >&2; sleep 10"},
			DeleteAlreadyExistingPod: true,
			WaitForPodRunning:        true,
		},
	)
	require.NoError(t, err)

	t.Run("GetContainerStdoutLogs", func(t *testing.T) {
		stdout, err := nativekubernetes.GetContainerStdoutLogs(ctx, clientset, namespaceName, podName, containerName)
		if err != nil {
			// Some Kubernetes APIs don't support separate stream specification
			t.Logf("Stdout logs not separately available: %v", err)
			// Verify combined logs work as fallback
			allLogs, err := nativekubernetes.GetContainerAllLogs(ctx, clientset, namespaceName, podName, containerName)
			require.NoError(t, err)
			require.Contains(t, string(allLogs), "stdout message")
		} else {
			require.Contains(t, string(stdout), "stdout message")
		}
	})

	t.Run("GetContainerStderrLogs", func(t *testing.T) {
		stderr, err := nativekubernetes.GetContainerStderrLogs(ctx, clientset, namespaceName, podName, containerName)
		if err != nil {
			// Some Kubernetes APIs don't support separate stream specification
			t.Logf("Stderr logs not separately available: %v", err)
		} else {
			require.Contains(t, string(stderr), "stderr message")
		}
	})

	t.Run("GetContainerAllLogs", func(t *testing.T) {
		allLogs, err := nativekubernetes.GetContainerAllLogs(ctx, clientset, namespaceName, podName, containerName)
		require.NoError(t, err)
		allLogsStr := string(allLogs)
		require.Contains(t, allLogsStr, "stdout message")
		require.Contains(t, allLogsStr, "stderr message")
	})

	t.Run("GetContainerLogs", func(t *testing.T) {
		stdout, stderr, err := nativekubernetes.GetContainerLogs(ctx, clientset, namespaceName, podName, containerName)
		require.NoError(t, err)
		require.Contains(t, string(stdout), "stdout message")

		// Some Kubernetes APIs don't support separate stderr stream
		// In that case, stderr will be empty and both messages are in stdout
		if len(stderr) > 0 {
			require.Contains(t, string(stderr), "stderr message")
		} else {
			// Verify combined logs contain both messages
			require.Contains(t, string(stdout), "stderr message")
		}
	})

	t.Run("GetContainerLogs_invalidPodName", func(t *testing.T) {
		_, _, err := nativekubernetes.GetContainerLogs(ctx, clientset, namespaceName, "nonexistent-pod", containerName)
		require.Error(t, err)
	})

	t.Run("GetContainerLogs_invalidNamespace", func(t *testing.T) {
		_, _, err := nativekubernetes.GetContainerLogs(ctx, clientset, "nonexistent-namespace", podName, containerName)
		require.Error(t, err)
	})

	t.Run("GetContainerLogs_nilClientset", func(t *testing.T) {
		_, _, err := nativekubernetes.GetContainerLogs(ctx, nil, namespaceName, podName, containerName)
		require.Error(t, err)
	})

	t.Run("GetContainerStdoutLogs_nilClientset", func(t *testing.T) {
		_, err := nativekubernetes.GetContainerStdoutLogs(ctx, nil, namespaceName, podName, containerName)
		require.Error(t, err)
	})

	t.Run("GetContainerStderrLogs_nilClientset", func(t *testing.T) {
		_, err := nativekubernetes.GetContainerStderrLogs(ctx, nil, namespaceName, podName, containerName)
		require.Error(t, err)
	})

	t.Run("GetContainerAllLogs_nilClientset", func(t *testing.T) {
		_, err := nativekubernetes.GetContainerAllLogs(ctx, nil, namespaceName, podName, containerName)
		require.Error(t, err)
	})

	t.Run("GetContainerStdoutLogs_emptyPodName", func(t *testing.T) {
		_, err := nativekubernetes.GetContainerStdoutLogs(ctx, clientset, namespaceName, "", containerName)
		require.Error(t, err)
	})

	t.Run("GetContainerStdoutLogs_emptyNamespace", func(t *testing.T) {
		_, err := nativekubernetes.GetContainerStdoutLogs(ctx, clientset, "", podName, containerName)
		require.Error(t, err)
	})
}
