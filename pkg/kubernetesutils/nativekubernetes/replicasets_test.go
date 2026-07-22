package nativekubernetes_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_CreateAndDeleteReplicaSet(t *testing.T) {
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

		const replicaSetName = "testreplicaset"
		const namespaceName = "default"

		err = nativekubernetes.DeleteReplicaSet(ctx, clientset, replicaSetName, namespaceName)
		require.NoError(t, err)

		exists, err := nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
		require.NoError(t, err)
		require.False(t, exists)

		// check if consecutive create, delete, create, delete... works
		for range 3 {
			err = nativekubernetes.CreateReplicaSet(ctx, config, &kubernetesparameteroptions.RunCommandOptions{
				Namespace:      namespaceName,
				ReplicaSetName: replicaSetName,
				Image:          "ubuntu",
				Command:        []string{"bash", "-c", "sleep 1m"},
				Replicas:       2,
			})
			require.NoError(t, err)

			exists, err = nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
			require.NoError(t, err)
			require.True(t, exists)

			err = nativekubernetes.DeleteReplicaSet(ctx, clientset, replicaSetName, namespaceName)
			require.NoError(t, err)

			exists, err = nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
			require.NoError(t, err)
			require.False(t, exists)
		}
	})
}

func Test_WaitForReplicaSetDeleted(t *testing.T) {
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

	t.Run("already deleted ReplicaSet", func(t *testing.T) {
		replicaSetName := "testreplicaset"
		namespaceName := "default"

		// Ensure ReplicaSet is absent
		err = nativekubernetes.DeleteReplicaSet(ctx, clientset, replicaSetName, namespaceName)
		require.NoError(t, err)

		// Check there's no wait for an already deleted ReplicaSet:
		err = nativekubernetes.WaitForReplicaSetDeleted(ctx, clientset, replicaSetName, namespaceName, time.Second*1)
		require.NoError(t, err)
	})
}
