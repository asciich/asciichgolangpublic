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
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

// This example shows how to run a single command in a temporary pod with a secret mounted.
// The secret is created first, then referenced as an environment variable in the pod.
// The pod is created automatically, executes the command, and is deleted afterwards.
func Test_Example_RunSingleCommandPodWithSecret(t *testing.T) {
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
	const podName = "example-command-pod-with-secret"
	const secretName = "example-secret"
	const secretKey = "mykey"
	const secretValue = "mysecretvalue"
	const envVarName = "MY_SECRET_VAR"
	const namespaceName = "default"

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

	// Run a single command in a temporary pod with the secret mounted as environment variable:
	output, err := nativekubernetes.RunCommandInTemporaryPod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			Image:                    "ubuntu",
			PodName:                  podName,
			DeleteAlreadyExistingPod: true,
			Command:                  []string{"bash", "-c", "echo Secret value is: $" + envVarName},
			SecretEnvVars: map[string]kubernetesparameteroptions.SecretEnvVarSource{
				envVarName: {
					SecretName: secretName,
					SecretKey:  secretKey,
				},
			},
		},
	)
	require.NoError(t, err)

	// Read the stdout of the executed command:
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.Contains(t, stdout, "Secret value is:")
	require.Contains(t, stdout, secretValue)
}
