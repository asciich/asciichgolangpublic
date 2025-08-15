package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateAndDeletePod(t *testing.T) {
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
			err = nativekubernetes.CreatePod(ctx, clientset, &kubernetesparameteroptions.RunCommandOptions{
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
