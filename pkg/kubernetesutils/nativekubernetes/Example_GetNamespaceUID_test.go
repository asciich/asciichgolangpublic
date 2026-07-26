package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/uuidutils"
)

func Test_Example_GetNamespaceUid(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.ContextVerbose()

	// -----
	// Prepare test environment start ...
	// Use the shared cluster for concurrent-safe test execution
	const clusterName = kindutils.SharedClusterName

	// Get or create the shared cluster (uses file-based locking)
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Create a unique namespace for this test to avoid conflicts
	namespaceName := "test-get-namespace-uid"
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	// Get the clientSet to access the kubernetes cluster:
	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Get the namespace UID:
	uid, err := nativekubernetes.GetNamespaceUid(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.True(t, uuidutils.IsUuid(uid))
}
