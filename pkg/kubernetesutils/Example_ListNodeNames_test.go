package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kindutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_ListNodeNames(t *testing.T) {
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

	// Get Kubernetes cluster:
	cluster, err := nativekubernetes.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// List all node names:
	nodeNames, err := cluster.ListNodeNames(ctx)
	require.NoError(t, err)

	// We expect the "default" namespace to be present:
	require.Contains(t, nodeNames, "kubernetesutils-control-plane")

}
