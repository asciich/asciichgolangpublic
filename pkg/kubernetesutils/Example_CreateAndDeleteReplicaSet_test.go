package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CreateAndDeleteReplicaSet(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := testutils.GetKindClusterNameForTest(t)

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Define the namespace and replicaSet name we use for testing
	const namespaceName = "testnamespace"
	const replicaSetName = "example-replicaset"

	// Get the namespace (will be created implicitly when creating the replicaSet)
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the replicaSet object
	replicaSet, err := namespace.GetReplicaSetByName(replicaSetName)
	require.NoError(t, err)

	// Ensure the replicaSet is absent
	err = replicaSet.Delete(ctx)
	require.NoError(t, err)
	exists, err := replicaSet.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example replicaSet using YAML
	replicaSetYaml := ""
	replicaSetYaml += "apiVersion: apps/v1\n"
	replicaSetYaml += "kind: ReplicaSet\n"
	replicaSetYaml += "metadata:\n"
	replicaSetYaml += "  name: " + replicaSetName + "\n"
	replicaSetYaml += "  namespace: " + namespaceName + "\n"
	replicaSetYaml += "spec:\n"
	replicaSetYaml += "  replicas: 2\n"
	replicaSetYaml += "  selector:\n"
	replicaSetYaml += "    matchLabels:\n"
	replicaSetYaml += "      app: " + replicaSetName + "\n"
	replicaSetYaml += "  template:\n"
	replicaSetYaml += "    metadata:\n"
	replicaSetYaml += "      labels:\n"
	replicaSetYaml += "        app: " + replicaSetName + "\n"
	replicaSetYaml += "    spec:\n"
	replicaSetYaml += "      containers:\n"
	replicaSetYaml += "      - name: ubuntu\n"
	replicaSetYaml += "        image: ubuntu\n"
	replicaSetYaml += "        command: [\"bash\", \"-c\", \"sleep 1m\"]\n"

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: replicaSetYaml,
	})
	require.NoError(t, err)
	exists, err = replicaSet.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the replicaSet exists via the cluster
	names, err := namespace.ListObjectNames(&kubernetesparameteroptions.ListKubernetesObjectsOptions{
		Namespace:  namespaceName,
		ObjectType: "replicaset",
	})
	require.NoError(t, err)
	require.Contains(t, names, replicaSetName)

	// Delete the replicaSet
	err = replicaSet.Delete(ctx)
	require.NoError(t, err)
	exists, err = replicaSet.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)
}
