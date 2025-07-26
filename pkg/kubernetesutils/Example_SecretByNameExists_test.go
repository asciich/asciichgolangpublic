package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_SecretByNameExists(t *testing.T) {
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
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Create an example secret. This implicitly generates the namespace if it does not exist.
	const namespaceName = "testnamespace"
	const secretName = "example-secret"
	_, err = cluster.CreateSecret(ctx, namespaceName, secretName, &kubernetesparameteroptions.CreateSecretOptions{
		SecretData: map[string][]byte{"my-secret": []byte("very-secret")},
	})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created secret exists:
	exists, err := cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.True(t, exists)

	// The same secret name in the default namespace does not exist:
	exists, err = cluster.SecretByNameExists(ctx, "default", secretName)
	require.NoError(t, err)
	require.False(t, exists)

	// This secret is expected to be in the same namespace but does not exist:
	exists, err = cluster.SecretByNameExists(ctx, namespaceName, "secret-does-not-exist")
	require.NoError(t, err)
	require.False(t, exists)

	// If we delete our secret again...
	err = cluster.DeleteSecretByName(ctx, namespaceName, secretName)
	require.NoError(t, err)

	// ... our secret becomes absent:
	exists, err = cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.False(t, exists)
}
