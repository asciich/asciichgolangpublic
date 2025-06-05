package kubernetesutils_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_WatchConfigMap(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	const clusterName = "kind"

	// Ensure a local kind cluster is available for testing:
	_, err := commandexecutor.Bash().RunOneLiner(ctx, fmt.Sprintf("kind create cluster -n '%s' || true", clusterName))
	require.NoError(t, err)

	// Get Kubernetes cluster:
	cluster, err := nativekubernetes.GetClusterByName(ctx, clusterName)
	require.NoError(t, err)

	// Ensure namespace exists
	const namespaceName = "testnamespace"
	namespace, err := cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Ensure ConfigMap under test is absent:
	const configmapName = "example-configmap"
	err = namespace.DeleteConfigMapByName(ctx, configmapName)
	require.NoError(t, err)

	// define counters and watch config map
	var cmCreateCounter, cmUpdateCounter, cmDeleteCounter int
	ctxWatch, cancel := context.WithCancel(ctx) // ensure we can cancel the watching
	err = namespace.WatchConfigMap(
		ctxWatch,
		configmapName,
		func(kubernetesutils.ConfigMap) { cmCreateCounter++ },
		func(kubernetesutils.ConfigMap) { cmUpdateCounter++ },
		func(kubernetesutils.ConfigMap) { cmDeleteCounter++ },
	)
	require.NoError(t, err)
	defer cancel()

	// check no callback called
	time.Sleep(100 * time.Millisecond)
	require.EqualValues(t, 0, cmCreateCounter)
	require.EqualValues(t, 0, cmUpdateCounter)
	require.EqualValues(t, 0, cmDeleteCounter)

	// create config map
	_, err = cluster.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesutils.CreateConfigMapOptions{
		ConfigMapData: map[string]string{"my-configmap": "configmap content"},
	})
	require.NoError(t, err)

	// Check cmCreateCounter is incremented
	time.Sleep(100 * time.Millisecond)
	require.EqualValues(t, 1, cmCreateCounter)
	require.EqualValues(t, 0, cmUpdateCounter)
	require.EqualValues(t, 0, cmDeleteCounter)

	var nUpdates = 3
	for i := 0; i < nUpdates; i++ {
		// update config map
		_, err = cluster.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesutils.CreateConfigMapOptions{
			ConfigMapData: map[string]string{"my-configmap": "configmap content" + strconv.Itoa(i)},
		})
		require.NoError(t, err)

		// Check cmUpdateCounter is incremented
		time.Sleep(100 * time.Millisecond)
		require.EqualValues(t, 1, cmCreateCounter)
		require.EqualValues(t, i+1, cmUpdateCounter)
		require.EqualValues(t, 0, cmDeleteCounter)
	}

	// delete config map
	err = cluster.DeleteConfigMapByName(ctx, namespaceName, configmapName)
	require.NoError(t, err)

	// Check cmDeleteCounter is incremented
	time.Sleep(100 * time.Millisecond)
	require.EqualValues(t, 1, cmCreateCounter)
	require.EqualValues(t, nUpdates, cmUpdateCounter)
	require.EqualValues(t, 1, cmDeleteCounter)

	// cancel watching
	cancel()

	// Do further updates
	for i := 0; i < nUpdates; i++ {
		// update config map
		_, err = cluster.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesutils.CreateConfigMapOptions{
			ConfigMapData: map[string]string{"my-configmap": "configmap content"},
		})
		require.NoError(t, err)

		// Check cmUpdateCounter is unchanged since watch was deactivated
		time.Sleep(100 * time.Millisecond)
		require.EqualValues(t, 1, cmCreateCounter)
		require.EqualValues(t, nUpdates, cmUpdateCounter)
		require.EqualValues(t, 1, cmDeleteCounter)
	}

	// delete config map again
	err = cluster.DeleteConfigMapByName(ctx, namespaceName, configmapName)
	require.NoError(t, err)

	// Check cmUpdateCounter is unchanged since watch was deactivated
	time.Sleep(100 * time.Millisecond)
	require.EqualValues(t, 1, cmCreateCounter)
	require.EqualValues(t, nUpdates, cmUpdateCounter)
	require.EqualValues(t, 1, cmDeleteCounter)
}
