package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// This example shows how to run a single command in a temporary pod with a secret mounted as a file.
// The secret is created first, then mounted to a specific path in the container.
// The pod is created automatically, executes the command, and is deleted afterwards.
func Test_Example_RunSingleCommandPodWithSecretAsFile(t *testing.T) {
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

	// Get Kubernetes clientset:
	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Define names
	const podName = "example-command-pod-with-secret-file"
	const secretName = "example-secret-file"
	const secretKey = "mykey"
	const secretValue = "mysecretvalue"
	const mountPath = "/etc/secret"
	const secretFilePath = mountPath + "/" + secretKey
	const namespaceName = "kubernetesutils-runwithsecret"

	// Create the secret first:
	_, err = cluster.CreateSecret(ctx, namespaceName, secretName, &kubernetesparameteroptions.CreateSecretOptions{
		SecretData: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	})
	require.NoError(t, err)

	// Ensure secret is deleted at the end:
	defer func() {
		_ = cluster.DeleteSecretByName(ctx, namespaceName, secretName)
	}()

	// Verify secret exists:
	exists, err := cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.True(t, exists)

	// Run a single command in a temporary pod with the secret mounted as a file:
	output, err := nativekubernetes.RunCommandInTemporaryPod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			SecretMounts: map[string]kubernetesparameteroptions.SecretMountSource{
				mountPath: {
					SecretName: secretName,
				},
			},
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "cat " + secretFilePath},
			},
		},
	)
	require.NoError(t, err)

	// Read the stdout of the executed command:
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.EqualValues(t, secretValue, stdout)
}
