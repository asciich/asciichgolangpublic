package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_CreateAndDeleteNamespace(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

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
	const namespaceName = "create-and-delete-ns"

	// Ensure the namespace is absent:
	err = nativekubernetes.DeleteNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	exists, err := nativekubernetes.NamespaceExists(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Create the namespace:
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.NamespaceExists(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// The create function is idempotent.
	// If the namespace does already exists no error will be returned:
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.NamespaceExists(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Delete the namespace again so we have no leftovers:
	err = nativekubernetes.DeleteNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.NamespaceExists(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// The delete function is idempotent as well.
	// If the namespace is already absent no error will be raised.
	err = nativekubernetes.DeleteNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.NamespaceExists(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)
}
