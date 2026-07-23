package nativekubernetes_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_CreateAndDeletePod(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("happy path", func(t *testing.T) {

		const podName = "testpod"
		const namespaceName = "default"

		err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
		require.NoError(t, err)

		exists, err := nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
		require.NoError(t, err)
		require.False(t, exists)

		// check if consecutive create, delete, create, delete... works
		for range 3 {
			err = nativekubernetes.CreatePod(ctx, config, &kubernetesparameteroptions.RunCommandOptions{
				Namespace: namespaceName,
				PodName:   podName,
				Image:     "ubunt",
				Command:   []string{"bash", "-c", "sleep 1m"},
			})
			require.NoError(t, err)

			exists, err = nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
			require.NoError(t, err)
			require.True(t, exists)

			err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
			require.NoError(t, err)

			exists, err = nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
			require.NoError(t, err)
			require.False(t, exists)
		}
	})
}

func Test_WaitForDeleted(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("already deleted pod", func(t *testing.T) {
		podName := "testpod"
		namespaceName := "default"

		// Ensure pod is absent
		err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
		require.NoError(t, err)

		// Check there's no wait for an already deleted pod:
		err = nativekubernetes.WaitForPodDeleted(ctx, clientset, podName, namespaceName, time.Second*1)
		require.NoError(t, err)
	})
}

func Test_ListPods(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("create and delete pods with list in between", func(t *testing.T) {
		const namespaceName = "default"

		podNames := []string{"listpod-1", "listpod-2", "listpod-3"}

		// Ensure all test pods are absent before starting
		for _, name := range podNames {
			err = nativekubernetes.DeletePod(ctx, clientset, name, namespaceName)
			require.NoError(t, err)
		}

		// Create pods one by one and verify list grows
		for i, name := range podNames {
			err = nativekubernetes.CreatePod(ctx, config, &kubernetesparameteroptions.RunCommandOptions{
				Namespace: namespaceName,
				PodName:   name,
				Image:     "ubuntu",
				Command:   []string{"bash", "-c", "sleep 1m"},
			})
			require.NoError(t, err)

			listed, err := nativekubernetes.ListPods(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, created := range podNames[:i+1] {
				require.Contains(t, listed, created)
			}
			for _, notYetCreated := range podNames[i+1:] {
				require.NotContains(t, listed, notYetCreated)
			}
		}

		// Delete pods one by one and verify list shrinks
		for i, name := range podNames {
			err = nativekubernetes.DeletePod(ctx, clientset, name, namespaceName)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListPods(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, deleted := range podNames[:i+1] {
				require.NotContains(t, listed, deleted)
			}
			for _, stillPresent := range podNames[i+1:] {
				require.Contains(t, listed, stillPresent)
			}
		}
	})
}
