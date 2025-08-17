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

// This example shows how to exec in a container and write to the stdin of this command
func Test_Example_WriteToStdinOfExecCommand(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// define the pod name and the namespace
	const podName = "example-run-pod"
	const namespaceName = "default"

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	// Get the config and clientset to access the kubernetes cluster:
	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Ensure the pod is absent before the test starts:
	err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)

	// Ensure pod is deleted at the end of the test:
	defer func() { nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName) }()

	// Check pod absent
	exists, err := nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// start the pod
	err = nativekubernetes.CreatePod(
		ctx,
		config,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			Namespace:                "default",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"}, // Same as sleep inifinity but does not ignore SIGINT
			WaitForPodRunning:        true,
		},
	)
	require.NoError(t, err)

	// Check pod present
	exists, err = nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Exec command printing out what's written to stdin
	output, err := nativekubernetes.Exec(
		ctx,
		config,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			Namespace:                "default",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"cat"},
			StdinBytes:               []byte("hello world"),
		},
	)
	require.NoError(t, err)

	// Check the output of the exec command:
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.EqualValues(t, "hello world", stdout)
}
