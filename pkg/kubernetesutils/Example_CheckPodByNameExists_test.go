package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CheckPodByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Create an example namespace and pod.
	const namespaceName = "testnamespace"
	const podName = "example-pod"
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	_, err = namespace.CreatePod(ctx, &kubernetesparameteroptions.RunCommandOptions{
		Image:   "busybox",
		Command: []string{"sleep", "3600"},
		PodName: podName,
	})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created pod exists - CheckPodByNameExists returns nil:
	err = cluster.CheckPodByNameExists(ctx, namespaceName, podName)
	require.NoError(t, err)

	// The same pod name in the default namespace does not exist - CheckPodByNameExists returns error:
	err = cluster.CheckPodByNameExists(ctx, "default", podName)
	require.Error(t, err)

	// This pod is expected to be in the same namespace but does not exist:
	err = cluster.CheckPodByNameExists(ctx, namespaceName, "pod-does-not-exist")
	require.Error(t, err)

	// If we delete our pod again...
	err = namespace.DeletePodByName(ctx, podName)
	require.NoError(t, err)

	// ... our pod becomes absent and CheckPodByNameExists returns error:
	err = cluster.CheckPodByNameExists(ctx, namespaceName, podName)
	require.Error(t, err)
}
