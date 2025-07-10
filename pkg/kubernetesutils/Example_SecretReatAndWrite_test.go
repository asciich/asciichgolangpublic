package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_ReadAndWriteSecret(t *testing.T) {
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
	cluster, err := nativekubernetes.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// define the namespace and secret name we use for testing
	const namespaceName = "testnamespace"
	const secretName = "example-secret"

	// Ensure the secret is absent
	err = cluster.DeleteSecretByName(ctx, namespaceName, secretName)
	require.NoError(t, err)
	exists, err := cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example secret. This implicitly generates the namespace if it does not exist.
	secret, err := cluster.CreateSecret(ctx, namespaceName, secretName, &kubernetesparameteroptions.CreateSecretOptions{
		SecretData: map[string][]byte{"my-secret": []byte("very-secret")},
	})
	require.NoError(t, err)
	exists, err = cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.True(t, exists)

	// Read the secret directly on the cluster:
	secretData, err := cluster.ReadSecret(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.Len(t, secretData, 1)
	require.EqualValues(t, secretData["my-secret"], []byte("very-secret"))

	// Or read the secret:
	secretData, err = secret.Read(ctx)
	require.NoError(t, err)
	require.Len(t, secretData, 1)
	require.EqualValues(t, secretData["my-secret"], []byte("very-secret"))

	// To update the secret we can simply call CreateSecret again:
	_, err = cluster.CreateSecret(ctx, namespaceName, secretName, &kubernetesparameteroptions.CreateSecretOptions{
		SecretData: map[string][]byte{"my-secret": []byte("very-secret2")},
	})
	require.NoError(t, err)
	exists, err = cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.True(t, exists)

	// Read the secret
	secretData, err = cluster.ReadSecret(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.Len(t, secretData, 1)
	require.EqualValues(t, secretData["my-secret"], []byte("very-secret2"))
}
