package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CheckNamespaceByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Create an example namespace.
	const namespaceName = "example-namespace"
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created namespace exists - CheckNamespaceByNameExists returns nil:
	err = cluster.CheckNamespaceByNameExists(ctx, namespaceName)
	require.NoError(t, err)

	// This namespace does not exist - CheckNamespaceByNameExists returns error:
	err = cluster.CheckNamespaceByNameExists(ctx, "namespace-does-not-exist")
	require.Error(t, err)

	// If we delete our namespace again...
	err = cluster.DeleteNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// ... our namespace becomes absent and CheckNamespaceByNameExists returns error:
	err = cluster.CheckNamespaceByNameExists(ctx, namespaceName)
	require.Error(t, err)
}
