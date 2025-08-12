package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

// This example shows how to run a pod and get the stdout of the executed command.
func Test_Example_RunPodAndGetStdout(t *testing.T) {
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

	// Get the config to access the kubernetes cluster:
	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// define the pod name
	const podName = "example-run-pod"

	// Run a command in a temporary pod:
	output, err := nativekubernetes.RunCommandInTemporaryPod(
		ctx,
		config,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			Namespace:                "default",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"bash", "-c", "echo hello_world"},
		},
	)
	require.NoError(t, err)

	// Read the stdout of the executed command:
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.EqualValues(t, "hello_world\n", stdout)
}
