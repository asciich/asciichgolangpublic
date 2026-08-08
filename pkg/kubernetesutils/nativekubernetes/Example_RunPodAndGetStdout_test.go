package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// This example shows how to run a pod and get the stdout of the executed command.
func Test_Example_RunPodAndGetStdout(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	namespaceName := "test-get-namespace-uid"

	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// define the pod name
	const podName = "example-run-pod"

	// Run a command in a temporary pod:
	output, err := nativekubernetes.RunCommandInTemporaryPod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "echo hello_world"},
			},
		},
	)
	require.NoError(t, err)

	// Read the stdout of the executed command:
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.EqualValues(t, "hello_world\n", stdout)
}
