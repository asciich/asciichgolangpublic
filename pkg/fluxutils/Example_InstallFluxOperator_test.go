package fluxutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

// This example shows how flux can be installed using the flux-operator.
//
// To start this test use:
//
//	bash -c "cd pkg/fluxutils && go test -v -run Test_InstallFluxOperator"
func Test_InstallFluxOperator(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := continuousintegration.GetDefaultKindClusterName()
	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)
	defer kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, clusterName)
	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Ensure flux is absent/ The namespace containg flux is deleted to showcase an installation:
	err = cluster.DeleteNamespaceByName(ctx, "flux-system")
	require.NoError(t, err)

	// Check if the "flux-system" namespace is absent:
	exists, err := cluster.NamespaceByNameExists(ctx, "flux-system")
	require.NoError(t, err)
	require.False(t, exists)

	// Install flux using flux-operator
	_, err = fluxutils.InstallFlux(ctx, &fluxparameteroptions.InstalFluxOptions{
		KubernetesCluster: cluster,
		Namespace:         "flux-system",
	})
	require.NoError(t, err)

	// Check if the "flux-system" namespace is prsent:
	exists, err = cluster.NamespaceByNameExists(ctx, "flux-system")
	require.NoError(t, err)
	require.True(t, exists)
}
