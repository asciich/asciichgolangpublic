package kubernetesutils_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

// Test_Example_WaitPodReady demonstrates how to wait for a pod to become ready.
// This example shows:
// 1. Creating a test pod
// 2. Waiting for the pod to reach Running phase
// 3. Cleaning up the pod after the test
func Test_Example_WaitPodReady(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Get Kubernetes cluster
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Define the namespace and pod name we use for testing
	const namespaceName = "wait-pod-test"
	const podName = "example-wait-pod"

	// Ensure the namespace exists
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// Get the namespace for pod operations
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Ensure the pod is absent before testing
	err = namespace.DeletePodByName(ctx, podName)
	require.NoError(t, err)
	exists, err := namespace.PodByNameExists(ctx, podName)
	require.NoError(t, err)
	require.False(t, exists)

	// Create a test pod that will run for 60 seconds
	podYaml := `
apiVersion: v1
kind: Pod
metadata:
  name: ` + podName + `
  namespace: ` + namespaceName + `
spec:
  containers:
  - name: test-container
    image: ubuntu
    command: ["bash", "-c", "sleep 60"]
`

	// Create the pod
	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: podYaml,
	})
	require.NoError(t, err, "Failed to create test pod")

	// Clean up pod after the test
	defer func() {
		err := namespace.DeletePodByName(ctx, podName)
		require.NoError(t, err)
	}()

	// Wait for the pod to become ready (Running phase)
	err = namespace.WaitUntilPodReady(ctx, podName, 60*time.Second)
	require.NoError(t, err, "Pod should become ready within timeout")

	// Verify the pod is actually in Running phase
	exists, err = namespace.PodByNameExists(ctx, podName)
	require.NoError(t, err)
	require.True(t, exists, "Pod should exist and be running")

	t.Logf("Pod '%s/%s' successfully reached Running phase", namespaceName, podName)
}
