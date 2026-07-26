package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_ListNamespaceNames(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// List all namespace names:
	namespaces, err := cluster.ListNamespaceNames(ctx)
	require.NoError(t, err)

	// We expect the "default" namespace to be present:
	require.Contains(t, namespaces, "default")

	// We expect the "kube-system" namespace to be present:
	require.Contains(t, namespaces, "kube-system")
}
