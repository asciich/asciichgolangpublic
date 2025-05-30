package kubernetesutils_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_ConfigMapByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Prepare start ...
	const clusterName = "kind"

	// Ensure a local kind cluster is available for testing:
	_, err := commandexecutor.Bash().RunOneLiner(ctx, fmt.Sprintf("kind create cluster -n '%s' || true", clusterName))
	require.NoError(t, err)

	// Get Kubernetes cluster:
	cluster, err := nativekubernetes.GetClusterByName(ctx, clusterName)
	require.NoError(t, err)

	// Create an example configmap. This implicitly generates the namespace if it does not exist.
	const namespaceName = "testnamespace"
	const configmapName = "example-configmap"
	_, err = cluster.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesutils.CreateConfigMapOptions{
		ConfigMapData: map[string]string{"my-configmap": "configmap content"},
	})
	require.NoError(t, err)
	// ... prepare finished.

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
