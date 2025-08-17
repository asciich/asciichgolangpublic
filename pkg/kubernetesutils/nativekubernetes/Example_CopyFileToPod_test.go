package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
)

// This example shows how to copy a local file to a container running in a k8s pod.
func Test_Example_CopyFileToPod(t *testing.T) {
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
		clientset,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			Namespace:                "default",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"bash", "-c", "sleep infinity"},
			WaitForPodRunning:        true,
		},
	)
	require.NoError(t, err)

	// Check pod present
	exists, err = nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// define random test content
	testContent, err := randomgenerator.GetRandomString(1024)
	require.NoError(t, err)

	// create a temporary file
	localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
	require.NoError(t, err)

	// Copy file to container
	destPath := "/testfile"
	containerName := podName
	err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
	require.NoError(t, err)

	// Check file with correct content was uploaded:
	commandOutput, err := nativekubernetes.Exec(ctx, config, &kubernetesparameteroptions.RunCommandOptions{
		Namespace: namespaceName,
		PodName:   podName,
		Command:   []string{"cat", destPath},
	})
	require.NoError(t, err)

	readContent, err := commandOutput.GetStdoutAsString()
	require.NoError(t, err)

	require.EqualValues(t, testContent, readContent)
}
