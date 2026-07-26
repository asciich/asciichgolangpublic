package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CreateAndDeleteConfigMap(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Define the namespace and configmap name we use for testing
	const namespaceName = "testnamespace"
	const configMapName = "example-configmap"

	// Ensure the configmap is absent
	err = cluster.DeleteConfigMapByName(ctx, namespaceName, configMapName)
	require.NoError(t, err)
	exists, err := cluster.ConfigMapByNameExists(ctx, namespaceName, configMapName)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example configmap. This implicitly generates the namespace if it does not exist.
	configMap, err := cluster.CreateConfigMap(ctx, namespaceName, configMapName, &kubernetesparameteroptions.CreateConfigMapOptions{
		ConfigMapData: map[string]string{"key1": "value1", "key2": "value2"},
	})
	require.NoError(t, err)
	exists, err = cluster.ConfigMapByNameExists(ctx, namespaceName, configMapName)
	require.NoError(t, err)
	require.True(t, exists)

	// Read the configmap data
	data, err := configMap.GetAllData(ctx)
	require.NoError(t, err)
	require.Len(t, data, 2)
	require.EqualValues(t, "value1", data["key1"])
	require.EqualValues(t, "value2", data["key2"])

	// Update the configmap by calling CreateConfigMap again
	_, err = cluster.CreateConfigMap(ctx, namespaceName, configMapName, &kubernetesparameteroptions.CreateConfigMapOptions{
		ConfigMapData: map[string]string{"key1": "updated-value1"},
	})
	require.NoError(t, err)

	// Read the updated configmap data
	data, err = configMap.GetAllData(ctx)
	require.NoError(t, err)
	require.Len(t, data, 1)
	require.EqualValues(t, "updated-value1", data["key1"])

	// Delete the configmap
	err = cluster.DeleteConfigMapByName(ctx, namespaceName, configMapName)
	require.NoError(t, err)
	exists, err = cluster.ConfigMapByNameExists(ctx, namespaceName, configMapName)
	require.NoError(t, err)
	require.False(t, exists)
}
