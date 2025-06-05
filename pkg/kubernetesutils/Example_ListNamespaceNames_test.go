package kubernetesutils_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_ListNamespaceNames(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	const clusterName = "kind"

	// Ensure a local kind cluster is available for testing:
	_, err := commandexecutor.Bash().RunOneLiner(ctx, fmt.Sprintf("kind create cluster -n '%s' || true", clusterName))
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetes.GetClusterByName(ctx, clusterName)
	require.NoError(t, err)

	// List all namespace names:
	namespaces, err := cluster.ListNamespaceNames(ctx)
	require.NoError(t, err)

	// We expect the "default" namespace to be present:
	require.Contains(t, namespaces, "default")

	// We expect the "kube-system" namespace to be present:
	require.Contains(t, namespaces, "kube-system")
}
