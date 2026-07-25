package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CreateAndDeletePod(t *testing.T) {
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
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Define the namespace and pod name we use for testing
	const namespaceName = "testnamespace"
	const podName = "example-pod"

	// Get the namespace (will be created implicitly when creating the pod)
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the pod object
	pod, err := namespace.GetPodByName(podName)
	require.NoError(t, err)

	// Ensure the pod is absent
	err = pod.Delete(ctx)
	require.NoError(t, err)
	exists, err := pod.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example pod using YAML
	podYaml := ""
	podYaml += "apiVersion: v1\n"
	podYaml += "kind: Pod\n"
	podYaml += "metadata:\n"
	podYaml += "  name: " + podName + "\n"
	podYaml += "  namespace: " + namespaceName + "\n"
	podYaml += "spec:\n"
	podYaml += "  containers:\n"
	podYaml += "  - name: ubuntu\n"
	podYaml += "    image: ubuntu\n"
	podYaml += "    command: [\"bash\", \"-c\", \"sleep 1m\"]\n"

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: podYaml,
	})
	require.NoError(t, err)
	exists, err = pod.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the pod exists via the cluster
	names, err := namespace.ListObjectNames(&kubernetesparameteroptions.ListKubernetesObjectsOptions{
		Namespace:  namespaceName,
		ObjectType: "pod",
	})
	require.NoError(t, err)
	require.Contains(t, names, podName)

	// Delete the pod
	err = pod.Delete(ctx)
	require.NoError(t, err)
	exists, err = pod.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)
}
