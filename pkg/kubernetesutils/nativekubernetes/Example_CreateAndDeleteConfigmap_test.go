package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

// This example shows how to create and delete a ConfigMap.
func Test_Example_CreateAndDeleteConfigMap(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	// Define the name of the ConfigMap we use for testing
	const configMapName = "example-configmap"
	const namespaceName = "default"

	configMapData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// Ensure the ConfigMap is absent:
	err = nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, configMapName)
	require.NoError(t, err)

	configMapNames, err := nativekubernetes.ListConfigMaps(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, configMapNames, configMapName)

	// Create the ConfigMap:
	err = nativekubernetes.CreateConfigMap(ctx, clientset, namespaceName, configMapName, configMapData, nil)
	require.NoError(t, err)

	configMapNames, err = nativekubernetes.ListConfigMaps(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.Contains(t, configMapNames, configMapName)

	// Read the ConfigMap and verify the data:
	data, err := nativekubernetes.GetConfigMapData(ctx, clientset, namespaceName, configMapName)
	require.NoError(t, err)
	require.Equal(t, configMapData, data)

	// The create function is idempotent if we handle AlreadyExists:
	exists, err := nativekubernetes.ConfigMapExists(ctx, clientset, namespaceName, configMapName)
	require.NoError(t, err)
	require.True(t, exists)

	// Delete the ConfigMap again so we have no leftovers:
	err = nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, configMapName)
	require.NoError(t, err)

	configMapNames, err = nativekubernetes.ListConfigMaps(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, configMapNames, configMapName)

	// The delete function is idempotent as well.
	// If the ConfigMap is already absent no error will be raised.
	err = nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, configMapName)
	require.NoError(t, err)

	configMapNames, err = nativekubernetes.ListConfigMaps(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, configMapNames, configMapName)
}
