package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/uuidutils"
)

func Test_Example_GetNamespaceUid(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.ContextVerbose()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	// Get the clientSet to access the kubernetes cluster:
	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Define the name of the namespace we use for testing
	const namespaceName = "get-ns-uid"

	// Create the namespace:
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	exists, err := nativekubernetes.NamespaceExists(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Get the namespace UID:
	uid, err := nativekubernetes.GetNamespaceUid(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.True(t, uuidutils.IsUuid(uid))
}
