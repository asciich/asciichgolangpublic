package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateAndDeleteConfigMap(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("happy path", func(t *testing.T) {

		const configMapName = "testconfigmap"
		const namespaceName = "default"

		configMapData := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		err = nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, configMapName)
		require.NoError(t, err)

		exists, err := nativekubernetes.ConfigMapExists(ctx, clientset, namespaceName, configMapName)
		require.NoError(t, err)
		require.False(t, exists)

		// check if consecutive create, delete, create, delete... works
		for range 3 {
			err = nativekubernetes.CreateConfigMap(ctx, clientset, namespaceName, configMapName, configMapData, nil)
			require.NoError(t, err)

			exists, err = nativekubernetes.ConfigMapExists(ctx, clientset, namespaceName, configMapName)
			require.NoError(t, err)
			require.True(t, exists)

			err = nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, configMapName)
			require.NoError(t, err)

			exists, err = nativekubernetes.ConfigMapExists(ctx, clientset, namespaceName, configMapName)
			require.NoError(t, err)
			require.False(t, exists)
		}
	})
}

func Test_ListConfigMaps(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("create and delete ConfigMaps with list in between", func(t *testing.T) {
		const namespaceName = "default"

		configMapNames := []string{"listcm-1", "listcm-2", "listcm-3"}
		configMapData := map[string]string{"key": "value"}

		// Ensure all test ConfigMaps are absent before starting
		for _, name := range configMapNames {
			err = nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, name)
			require.NoError(t, err)
		}

		// Create ConfigMaps one by one and verify list grows
		for i, name := range configMapNames {
			err = nativekubernetes.CreateConfigMap(ctx, clientset, namespaceName, name, configMapData, nil)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListConfigMaps(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, created := range configMapNames[:i+1] {
				require.Contains(t, listed, created)
			}
			for _, notYetCreated := range configMapNames[i+1:] {
				require.NotContains(t, listed, notYetCreated)
			}
		}

		// Delete ConfigMaps one by one and verify list shrinks
		for i, name := range configMapNames {
			err = nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, name)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListConfigMaps(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, deleted := range configMapNames[:i+1] {
				require.NotContains(t, listed, deleted)
			}
			for _, stillPresent := range configMapNames[i+1:] {
				require.Contains(t, listed, stillPresent)
			}
		}
	})
}
