package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

// This example shows how to create and delete a ReplicaSet.
func Test_Example_CreateAndDeleteReplicaSet(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

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

	// Define the name of the ReplicaSet we use for testing
	const replicaSetName = "example-replicaset"
	const namespaceName = "default"

	// Ensure the ReplicaSet is absent:
	err = nativekubernetes.DeleteReplicaSet(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)

	exists, err := nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Create the ReplicaSet with 2 replicas:
	err = nativekubernetes.CreateReplicaSet(
		ctx,
		config,
		&kubernetesparameteroptions.RunCommandOptions{
			Namespace:      namespaceName,
			ReplicaSetName: replicaSetName,
			Image:          "ubuntu",
			Command:        []string{"bash", "-c", "sleep 1m"},
			Replicas:       2,
		},
	)
	require.NoError(t, err)

	exists, err = nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// The create function is idempotent if DeleteAlreadyExistingReplicaSet is set:
	err = nativekubernetes.CreateReplicaSet(
		ctx,
		config,
		&kubernetesparameteroptions.RunCommandOptions{
			Namespace:                       namespaceName,
			ReplicaSetName:                  replicaSetName,
			Image:                           "ubuntu",
			Command:                         []string{"bash", "-c", "sleep 1m"},
			Replicas:                        2,
			DeleteAlreadyExistingReplicaSet: true,
		},
	)
	require.NoError(t, err)

	exists, err = nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Delete the ReplicaSet again so we have no leftovers:
	err = nativekubernetes.DeleteReplicaSet(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// The delete function is idempotent as well.
	// If the ReplicaSet is already absent no error will be raised.
	err = nativekubernetes.DeleteReplicaSet(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.ReplicaSetExists(ctx, clientset, replicaSetName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)
}
