package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

// This example shows how to run a single command in a temporary pod.
// The pod is created automatically, executes the command, and is deleted afterwards.
func Test_Example_RunSingleCommandPod(t *testing.T) {
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

	// Get Kubernetes clientset:
	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Define the pod name and namespace
	const podName = "example-command-pod"
	const namespaceName = "default"

	// Run a single command in a temporary pod:
	output, err := nativekubernetes.RunCommandInTemporaryPod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"bash", "-c", "echo 'Hello from temporary pod'"},
		},
	)
	require.NoError(t, err)

	// Read the stdout of the executed command:
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.EqualValues(t, "Hello from temporary pod\n", stdout)
}
