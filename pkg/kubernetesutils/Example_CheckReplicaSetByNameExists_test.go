package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func Test_Example_CheckReplicaSetByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Create an example namespace and replicaSet.
	const namespaceName = "testnamespace"
	const replicaSetName = "example-replicaset"
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	_, err = namespace.CreateReplicaSet(ctx, &kubernetesparameteroptions.KubernetesRunCommandOptions{
		Image:          "busybox",
		ReplicaSetName: replicaSetName,
		Replicas:       1,
		RunCommandOptions: &parameteroptions.RunCommandOptions{
			Command: []string{"sleep", "3600"},
		},
	})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created replicaSet exists - CheckReplicaSetByNameExists returns nil:
	err = cluster.CheckReplicaSetByNameExists(ctx, namespaceName, replicaSetName)
	require.NoError(t, err)

	// The same replicaSet name in the default namespace does not exist - CheckReplicaSetByNameExists returns error:
	err = cluster.CheckReplicaSetByNameExists(ctx, "default", replicaSetName)
	require.Error(t, err)

	// This replicaSet is expected to be in the same namespace but does not exist:
	err = cluster.CheckReplicaSetByNameExists(ctx, namespaceName, "replicaset-does-not-exist")
	require.Error(t, err)

	// If we delete our replicaSet again...
	err = namespace.DeleteReplicaSetByName(ctx, replicaSetName)
	require.NoError(t, err)

	// ... our replicaSet becomes absent and CheckReplicaSetByNameExists returns error:
	err = cluster.CheckReplicaSetByNameExists(ctx, namespaceName, replicaSetName)
	require.Error(t, err)
}
