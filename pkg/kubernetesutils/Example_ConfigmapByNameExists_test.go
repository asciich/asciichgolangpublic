package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_ConfigMapByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := continuousintegration.GetDefaultKindClusterName()

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)
	defer kindutils.DeleteClusterByName(ctx, clusterName)
		// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetes.GetClusterByName(ctx, "kind-" + clusterName)
	require.NoError(t, err)

	// Create an example configmap. This implicitly generates the namespace if it does not exist.
	const namespaceName = "testnamespace"
	const configmapName = "example-configmap"
	_, err = cluster.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{
		ConfigMapData: map[string]string{"my-configmap": "configmap content"},
	})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created configmap exists:
	exists, err := cluster.ConfigMapByNameExists(ctx, namespaceName, configmapName)
	require.NoError(t, err)
	require.True(t, exists)

	// The same configmap name in the default namespace does not exist:
	exists, err = cluster.ConfigMapByNameExists(ctx, "default", configmapName)
	require.NoError(t, err)
	require.False(t, exists)

	// This configmap is expected to be in the same namespace but does not exist:
	exists, err = cluster.ConfigMapByNameExists(ctx, namespaceName, "configmap-does-not-exist")
	require.NoError(t, err)
	require.False(t, exists)

	// If we delete our configmap again...
	err = cluster.DeleteConfigMapByName(ctx, namespaceName, configmapName)
	require.NoError(t, err)

	// ... our configmap becomes absent:
	exists, err = cluster.ConfigMapByNameExists(ctx, namespaceName, configmapName)
	require.NoError(t, err)
	require.False(t, exists)
}
