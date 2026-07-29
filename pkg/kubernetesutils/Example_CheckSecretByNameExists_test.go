package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CheckSecretByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
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

	// Our created secret exists - CheckSecretByNameExists returns nil:
	err = cluster.CheckSecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)

	// The same secret name in the default namespace does not exist - CheckSecretByNameExists returns error:
	err = cluster.CheckSecretByNameExists(ctx, "default", secretName)
	require.Error(t, err)

	// This secret is expected to be in the same namespace but does not exist:
	err = cluster.CheckSecretByNameExists(ctx, namespaceName, "secret-does-not-exist")
	require.Error(t, err)

	// If we delete our secret again...
	err = cluster.DeleteSecretByName(ctx, namespaceName, secretName)
	require.NoError(t, err)

	// ... our secret becomes absent and CheckSecretByNameExists returns error:
	err = cluster.CheckSecretByNameExists(ctx, namespaceName, secretName)
	require.Error(t, err)
}
