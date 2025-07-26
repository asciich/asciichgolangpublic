package helmutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils"
	"github.com/asciich/asciichgolangpublic/pkg/helmutils/helmparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

// Example how to install a hemlchart.
// For this example the flux operator is installed.
//
//	Source: https://fluxcd.io/flux/installation/#install-the-flux-operator
//
// To run this test use:
//
//	bash -c "cd pkg/helmutils && go test -v -run Test_InstallHelmchart_FluxOperator"
func Test_InstallHelmchart_FluxOperator(t *testing.T) {
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

	// Ensure flux/ the namespace which contains flux is absent.
	err = cluster.DeleteNamespaceByName(ctx, "flux-system")
	require.NoError(t, err)

	// Check namespace is absent
	exists, err := cluster.NamespaceByNameExists(ctx, "flux-system")
	require.NoError(t, err)
	require.False(t, exists)

	// The install is idempotent and can be repeated multiple times.
	// Behind the scenes it runs a helm upgrade if needed.
	//
	// We run the install twice to show that it's idempotent.
	for i := 0; i < 2; i++ {
		// Deploy flux operator using helm
		// Equvalent helm command:
		// helm install flux-operator oci://ghcr.io/controlplaneio-fluxcd/charts/flux-operator --namespace flux-system --create-namespace
		err = helmutils.InstallHelmChart(ctx, &helmparameteroptions.InstallHelmChartOptions{
			KubernetesCluster: cluster,
			ChartReference:    "flux-operator",
			ChartUri:          "oci://ghcr.io/controlplaneio-fluxcd/charts/flux-operator",
			Namespace:         "flux-system",
		})
		require.NoError(t, err)

		// Check namespace was created
		exists, err = cluster.NamespaceByNameExists(ctx, "flux-system")
		require.NoError(t, err)
		require.True(t, exists)
	}
}
